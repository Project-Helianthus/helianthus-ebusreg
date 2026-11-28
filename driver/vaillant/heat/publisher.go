// Package heat publishes one thermal graph for a controller, or for a heat
// generator when no controller is present. The graph is resources and
// relations. Feeds and writes belong to later slices.
package heat

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"sync"

	semreg "github.com/Project-Helianthus/helianthus-semreg/semreg/v1"
	semdriver "github.com/Project-Helianthus/helianthus-semreg/semreg/v1/driver"
)

// Publisher owns one thermal batch sequence. Publish is safe for concurrent use.
type Publisher struct {
	mu   sync.Mutex
	next uint64
}

// Publish stages the statement's resources, relations and identity links on
// the session's owner batch. Pending operator links replace only the link of
// their binding. An empty pending set leaves that batch unchanged. A statement
// with no root identity publishes nothing.
func (p *Publisher) Publish(ctx context.Context, session semdriver.Session, statement Statement) error {
	if p == nil {
		return errors.New("heat: nil publisher")
	}
	if ctx == nil {
		return errors.New("heat: nil context")
	}
	if session == nil || session.Publisher() == nil {
		return errors.New("heat: nil session publisher")
	}
	if !chooseRoot(statement).ok {
		return nil
	}
	tx, err := session.Publisher().Begin(ctx)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Abort(ctx)
		}
	}()
	var buildErr error
	batch, err := session.Prepare(func(batch *semreg.PublicationBatch) {
		if batch.Sequence == "" {
			batch.Sequence = p.nextSequence()
		}
		source := batch.SourceID
		epoch := batch.SourceEpochID
		_, index, ok := buildResources(statement, originStamp{source: &source, epoch: &epoch})
		if !ok {
			return
		}
		binding := semreg.NativeBindingID("binding.thermal." + digestParts("binding", string(index.root)))
		stamp := originStamp{source: &source, epoch: &epoch, binding: &binding}
		resources, index, _ := buildResources(statement, stamp)
		relations, err := buildRelations(statement, index, stamp)
		if err != nil {
			buildErr = err
			return
		}
		if len(resources) > 0 {
			batch.ResourceUpserts = resources
		}
		if len(relations) > 0 {
			batch.RelationUpserts = relations
		}
		batch.BindingUpserts = []semreg.NativeBinding{{
			BindingID: binding, AssetID: batch.AssetID, SourceID: batch.SourceID,
			SourceEpochID: batch.SourceEpochID, DriverGeneration: batch.DriverGeneration,
			NativeResource: evidenceRef(index.rootEvidence), State: semreg.BindingCurrent, Revision: "1",
		}}
		root := index.root
		links := []semreg.IdentityLink{{
			AssetID: batch.AssetID, BindingID: binding, State: semreg.LinkCandidate,
			Basis: []semreg.EvidenceRef{evidenceRef(index.rootEvidence)}, Revision: "1", Resource: &root,
		}}
		sort.Slice(links, func(i, j int) bool { return links[i].BindingID < links[j].BindingID })
		batch.IdentityLinkUpserts = carryOperatorLinks(links, tx.OperatorLinks(batch.BindingUpserts))
	})
	if err != nil {
		return err
	}
	if buildErr != nil {
		return buildErr
	}
	if batch.BatchDigest == "" {
		digest, err := batch.ComputedDigest()
		if err != nil {
			return err
		}
		batch.BatchDigest = digest
	}
	if err := tx.Stage(batch); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	committed = true
	return nil
}

func (p *Publisher) nextSequence() semreg.Uint64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.next++
	return semreg.Uint64(strconv.FormatUint(p.next, 10))
}

// carryOperatorLinks replaces the built link of each pending binding and
// keeps every other link. An empty pending set returns the built slice.
func carryOperatorLinks(built, pending []semreg.IdentityLink) []semreg.IdentityLink {
	if len(pending) == 0 {
		return built
	}
	replaced := make(map[semreg.NativeBindingID]bool, len(pending))
	for _, link := range pending {
		replaced[link.BindingID] = true
	}
	out := make([]semreg.IdentityLink, 0, len(built)+len(pending))
	for _, link := range built {
		if !replaced[link.BindingID] {
			out = append(out, link)
		}
	}
	out = append(out, pending...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].BindingID < out[j].BindingID })
	return out
}

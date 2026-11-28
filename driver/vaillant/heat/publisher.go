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

// Publisher owns one thermal batch sequence and the bindings whose committed
// link is an operator link. Publish is safe for concurrent use.
type Publisher struct {
	mu       sync.Mutex
	next     uint64
	operator map[semreg.NativeBindingID]struct{}
}

// operatorLinkRefusal is the commit error when a batch link replaces a
// published operator link that the pending set does not carry.
const operatorLinkRefusal = "identity: operator identity links differ from the pending declarations"

func operatorLinkRejected(err error) bool {
	return err != nil && err.Error() == operatorLinkRefusal
}

// Publish stages the statement's resources, relations and identity links on
// the session's owner batch. Pending operator links replace only the link of
// their binding. An empty pending set leaves that batch unchanged, except a
// binding whose committed link is already the operator link: its candidate is
// left off the batch so the operator link stays. A statement with no root
// identity publishes nothing.
func (p *Publisher) Publish(ctx context.Context, session semdriver.Session, statement Statement) error {
	if p == nil {
		return errors.New("heat: nil publisher")
	}
	replaced, err := p.publish(ctx, session, statement)
	if err == nil || !operatorLinkRejected(err) || !replaced {
		return err
	}
	_, err = p.publish(ctx, session, statement)
	return err
}

func (p *Publisher) publish(ctx context.Context, session semdriver.Session, statement Statement) (bool, error) {
	if ctx == nil {
		return false, errors.New("heat: nil context")
	}
	if session == nil || session.Publisher() == nil {
		return false, errors.New("heat: nil session publisher")
	}
	if !chooseRoot(statement).ok {
		return false, nil
	}
	tx, err := session.Publisher().Begin(ctx)
	if err != nil {
		return false, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Abort(ctx)
		}
	}()
	var buildErr error
	var pending []semreg.IdentityLink
	var upserts []semreg.IdentityLink
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
		pending = append([]semreg.IdentityLink(nil), tx.OperatorLinks(batch.BindingUpserts)...)
		carried := omitCommittedOperatorLinks(carryOperatorLinks(links, pending), pending, p.committedOperatorLinks())
		batch.IdentityLinkUpserts = carried
		upserts = append([]semreg.IdentityLink(nil), carried...)
	})
	if err != nil {
		return false, err
	}
	if buildErr != nil {
		return false, buildErr
	}
	if batch.BatchDigest == "" {
		digest, err := batch.ComputedDigest()
		if err != nil {
			return false, err
		}
		batch.BatchDigest = digest
	}
	if err := tx.Stage(batch); err != nil {
		return false, err
	}
	if err := tx.Commit(ctx); err != nil {
		if operatorLinkRejected(err) {
			return p.noteReplacements(upserts, pending), err
		}
		return false, err
	}
	p.noteOperatorLinks(pending)
	committed = true
	return false, nil
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

// omitCommittedOperatorLinks drops a built link whose binding already has a
// committed operator link. A pending link for that binding, including a
// withdrawal, stays. An empty committed set returns links unchanged.
func omitCommittedOperatorLinks(links, pending []semreg.IdentityLink, committed map[semreg.NativeBindingID]struct{}) []semreg.IdentityLink {
	if len(committed) == 0 {
		return links
	}
	pendingID := make(map[semreg.NativeBindingID]bool, len(pending))
	for _, link := range pending {
		pendingID[link.BindingID] = true
	}
	out := make([]semreg.IdentityLink, 0, len(links))
	for _, link := range links {
		if _, held := committed[link.BindingID]; held && !pendingID[link.BindingID] {
			continue
		}
		out = append(out, link)
	}
	return out
}

func (p *Publisher) committedOperatorLinks() map[semreg.NativeBindingID]struct{} {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.operator) == 0 {
		return nil
	}
	out := make(map[semreg.NativeBindingID]struct{}, len(p.operator))
	for id := range p.operator {
		out[id] = struct{}{}
	}
	return out
}

func (p *Publisher) noteOperatorLinks(pending []semreg.IdentityLink) {
	if len(pending) == 0 {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.operator == nil {
		p.operator = map[semreg.NativeBindingID]struct{}{}
	}
	for _, link := range pending {
		if link.State == semreg.LinkWithdrawn {
			delete(p.operator, link.BindingID)
			continue
		}
		p.operator[link.BindingID] = struct{}{}
	}
}

func (p *Publisher) noteReplacements(upserts, pending []semreg.IdentityLink) bool {
	pendingID := make(map[semreg.NativeBindingID]bool, len(pending))
	for _, link := range pending {
		pendingID[link.BindingID] = true
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.operator == nil {
		p.operator = map[semreg.NativeBindingID]struct{}{}
	}
	noted := false
	for _, link := range upserts {
		if pendingID[link.BindingID] {
			continue
		}
		p.operator[link.BindingID] = struct{}{}
		noted = true
	}
	return noted
}

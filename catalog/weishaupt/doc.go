// Package weishaupt catalogues the Weishaupt WRSol solar-controller messages
// and the Kromschroeder/Weishaupt burner-management memory access as the
// public documentation states them. The WRSol identification states
// manufacturer byte 0x10; the memory access uses primary byte 0x50, which the
// Kromschroeder page states as its manufacturer byte 0x50 (and 0xC5 for the
// Weishaupt side).
//
// Every documented write or command has an encoder; the plane methods in
// providers/weishaupt call them by hand, and nothing is sent unless a method is
// called. The 0902/0903 secondary messages are passive in the cited TypeSpec
// source and are not catalogued here: the page gives their register count but
// no payload layout. The package decides no support and no qualification;
// undocumented replies are kept raw.
package weishaupt

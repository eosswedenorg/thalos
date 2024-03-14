package server

import "github.com/eoscanada/eos-go/ship"

// convert a ActionTrace to ActionTraceV1
func toActionTraceV1(trace *ship.ActionTrace) *ship.ActionTraceV1 {
	if trace_v0, ok := trace.Impl.(*ship.ActionTraceV0); ok {
		// convert to v1
		return &ship.ActionTraceV1{
			ActionOrdinal:        trace_v0.ActionOrdinal,
			CreatorActionOrdinal: trace_v0.CreatorActionOrdinal,
			Receipt:              trace_v0.Receipt,
			Receiver:             trace_v0.Receiver,
			Act:                  trace_v0.Act,
			ContextFree:          trace_v0.ContextFree,
			Elapsed:              trace_v0.Elapsed,
			Console:              trace_v0.Console,
			AccountRamDeltas:     trace_v0.AccountRamDeltas,
			Except:               trace_v0.Except,
			ErrorCode:            trace_v0.ErrorCode,
			ReturnValue:          []byte{},
		}
	}
	return trace.Impl.(*ship.ActionTraceV1)
}

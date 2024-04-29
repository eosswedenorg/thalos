package server

import (
	"github.com/pnx/antelope-go/ship"
)

// convert a ActionTrace to ActionTraceV1
func toActionTraceV1(trace *ship.ActionTrace) *ship.ActionTraceV1 {
	if trace.V0 != nil {
		// convert to v1
		return &ship.ActionTraceV1{
			ActionOrdinal:        trace.V0.ActionOrdinal,
			CreatorActionOrdinal: trace.V0.CreatorActionOrdinal,
			Receipt:              trace.V0.Receipt,
			Receiver:             trace.V0.Receiver,
			Act:                  trace.V0.Act,
			ContextFree:          trace.V0.ContextFree,
			Elapsed:              trace.V0.Elapsed,
			Console:              trace.V0.Console,
			AccountRamDeltas:     trace.V0.AccountRamDeltas,
			Except:               trace.V0.Except,
			ErrorCode:            trace.V0.ErrorCode,
			ReturnValue:          []byte{},
		}
	}
	return trace.V1
}

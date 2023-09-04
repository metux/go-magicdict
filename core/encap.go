package core

import (
    "github.com/metux/go-magicdict/api"
)

func encap(v api.Any, parent api.Entry) (api.Entry, error, bool) {
    switch val := v.(type) {
        // raw types
        case nil:
            return nil, nil, false
        case api.AnyMap:
            return NewDict(&val), nil, true
        case api.AnyList:
            return NewList(val), nil, true
        case int:
            return NewScalarInt(val), nil, false
        case float64:
            return NewScalarFloat(val), nil, false
        case string:
            return NewScalarStr(val), nil, false

        // entry types
        case api.Entry:
            return val, nil, false
    }
    return nil, api.ErrUnknownEntryType, false
}

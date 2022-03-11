
package abi_cache

import (
    "time"
    "strings"
    "github.com/go-redis/redis/v8"
    redis_cache "github.com/go-redis/cache/v8"
    eos "github.com/eoscanada/eos-go"

    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

var abiString = `
{
	"version": "eosio::abi/1.0",
	"types": [{
		"new_type_name": "new_type_name_1",
		"type": "name"
	}],
	"structs": [
	{
		"name": "struct_name_1",
		"base": "struct_name_2",
		"fields": [
			{"name":"struct_1_field_1", "type":"new_type_name_1"},
			{"name":"struct_1_field_2", "type":"struct_name_3"},
			{"name":"struct_1_field_3", "type":"string?"},
			{"name":"struct_1_field_4", "type":"string?"},
			{"name":"struct_1_field_5", "type":"struct_name_4[]"}
		]
   },{
		"name": "struct_name_2",
		"base": "",
		"fields": [
			{"name":"struct_2_field_1", "type":"string"}
		]
   },{
		"name": "struct_name_3",
		"base": "",
		"fields": [
			{"name":"struct_3_field_1", "type":"string"}
		]
   },{
		"name": "struct_name_4",
		"base": "",
		"fields": [
			{"name":"struct_4_field_1", "type":"string"}
		]
   }
	],
  "actions": [{
		"name": "action_name_1",
		"type": "struct_name_1",
		"ricardian_contract": ""
  }],
  "tables": [{
      "name": "table_name_1",
      "index_type": "i64",
      "key_names": [
        "key_name_1",
        "key_name_2"
      ],
      "key_types": [
        "string",
        "int"
      ],
      "type": "struct_name_1"
    }
  ]
}
`

func TestGetSet(t *testing.T) {

    c := New("abi.cache.test", &redis_cache.Options{
        Redis: redis.NewClient(&redis.Options{}),
         // Cache 10k keys for 1 minute.
        LocalCache: redis_cache.NewTinyLFU(10000, time.Minute),
    })

    abi, err := eos.NewABI(strings.NewReader(abiString))
    if err != nil {
        t.Error("Failed to build ABI", err)
    }

    err = c.Set("testaccount", abi, time.Minute)
    if err != nil {
        t.Error("Failed to set cache item", err)
    }

    c_abi, err := c.Get("testaccount")
    if err != nil {
        t.Error("Failed to get cache item", err)
    }

    assert.Equal(t, c_abi.Version, "eosio::abi/1.0")

    // Types
    assert.Equal(t, c_abi.Types[0].NewTypeName, "new_type_name_1")
    assert.Equal(t, c_abi.Types[0].Type, "name")

    // Structs
    assert.Equal(t, c_abi.Structs[0].Name, "struct_name_1")
    assert.Equal(t, c_abi.Structs[0].Base, "struct_name_2")
    assert.Equal(t, c_abi.Structs[0].Fields[0].Name, "struct_1_field_1")
    assert.Equal(t, c_abi.Structs[0].Fields[0].Type, "new_type_name_1")
    assert.Equal(t, c_abi.Structs[0].Fields[1].Name, "struct_1_field_2")
    assert.Equal(t, c_abi.Structs[0].Fields[1].Type, "struct_name_3")
    assert.Equal(t, c_abi.Structs[0].Fields[2].Name, "struct_1_field_3")
    assert.Equal(t, c_abi.Structs[0].Fields[2].Type, "string?")
    assert.Equal(t, c_abi.Structs[0].Fields[3].Name, "struct_1_field_4")
    assert.Equal(t, c_abi.Structs[0].Fields[3].Type, "string?")
    assert.Equal(t, c_abi.Structs[0].Fields[4].Name, "struct_1_field_5")
    assert.Equal(t, c_abi.Structs[0].Fields[4].Type, "struct_name_4[]")

    assert.Equal(t, c_abi.Structs[1].Name, "struct_name_2")
    assert.Equal(t, c_abi.Structs[1].Base, "")
    assert.Equal(t, c_abi.Structs[1].Fields[0].Name, "struct_2_field_1")
    assert.Equal(t, c_abi.Structs[1].Fields[0].Type, "string")

    assert.Equal(t, c_abi.Structs[2].Name, "struct_name_3")
    assert.Equal(t, c_abi.Structs[2].Base, "")
    assert.Equal(t, c_abi.Structs[2].Fields[0].Name, "struct_3_field_1")
    assert.Equal(t, c_abi.Structs[2].Fields[0].Type, "string")

    assert.Equal(t, c_abi.Structs[3].Name, "struct_name_4")
    assert.Equal(t, c_abi.Structs[3].Base, "")
    assert.Equal(t, c_abi.Structs[3].Fields[0].Name, "struct_4_field_1")
    assert.Equal(t, c_abi.Structs[3].Fields[0].Type, "string")

    // Actions
    assert.Equal(t, c_abi.Actions[0].Name, eos.ActN("action_name_1"))
    assert.Equal(t, c_abi.Actions[0].Type, "struct_name_1")
    assert.Equal(t, c_abi.Actions[0].RicardianContract, "")

    // Tables
    assert.Equal(t, c_abi.Tables[0].Name, eos.TableName("table_name_1"))
    assert.Equal(t, c_abi.Tables[0].Type, "struct_name_1")
    assert.Equal(t, c_abi.Tables[0].IndexType, "i64")
    assert.Equal(t, c_abi.Tables[0].KeyNames[0], "key_name_1")
    assert.Equal(t, c_abi.Tables[0].KeyNames[1], "key_name_2")
    assert.Equal(t, c_abi.Tables[0].KeyTypes[0], "string")
    assert.Equal(t, c_abi.Tables[0].KeyTypes[1], "int")
}

func TestCacheMiss(t *testing.T) {

    c := New("abi.cache.test", &redis_cache.Options{
        Redis: redis.NewClient(&redis.Options{}),
        // Cache 10k keys for 1 minute.
        LocalCache: redis_cache.NewTinyLFU(10000, time.Minute),
    })

    _, err := c.Get("nonexist")
    require.Error(t, err)
}

package abi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pnx/antelope-go/api"
	"github.com/pnx/antelope-go/chain"

	"github.com/eosswedenorg/thalos/internal/cache"
	"github.com/stretchr/testify/assert"
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

func assert_abi(t *testing.T, abi *chain.Abi) {
	assert.Equal(t, abi.Version, "eosio::abi/1.0")

	// Types
	assert.Equal(t, abi.Types[0].NewTypeName, "new_type_name_1")
	assert.Equal(t, abi.Types[0].Type, "name")

	// Structs
	assert.Equal(t, abi.Structs[0].Name, "struct_name_1")
	assert.Equal(t, abi.Structs[0].Base, "struct_name_2")
	assert.Equal(t, abi.Structs[0].Fields[0].Name, "struct_1_field_1")
	assert.Equal(t, abi.Structs[0].Fields[0].Type, "new_type_name_1")
	assert.Equal(t, abi.Structs[0].Fields[1].Name, "struct_1_field_2")
	assert.Equal(t, abi.Structs[0].Fields[1].Type, "struct_name_3")
	assert.Equal(t, abi.Structs[0].Fields[2].Name, "struct_1_field_3")
	assert.Equal(t, abi.Structs[0].Fields[2].Type, "string?")
	assert.Equal(t, abi.Structs[0].Fields[3].Name, "struct_1_field_4")
	assert.Equal(t, abi.Structs[0].Fields[3].Type, "string?")
	assert.Equal(t, abi.Structs[0].Fields[4].Name, "struct_1_field_5")
	assert.Equal(t, abi.Structs[0].Fields[4].Type, "struct_name_4[]")

	assert.Equal(t, abi.Structs[1].Name, "struct_name_2")
	assert.Equal(t, abi.Structs[1].Base, "")
	assert.Equal(t, abi.Structs[1].Fields[0].Name, "struct_2_field_1")
	assert.Equal(t, abi.Structs[1].Fields[0].Type, "string")

	assert.Equal(t, abi.Structs[2].Name, "struct_name_3")
	assert.Equal(t, abi.Structs[2].Base, "")
	assert.Equal(t, abi.Structs[2].Fields[0].Name, "struct_3_field_1")
	assert.Equal(t, abi.Structs[2].Fields[0].Type, "string")

	assert.Equal(t, abi.Structs[3].Name, "struct_name_4")
	assert.Equal(t, abi.Structs[3].Base, "")
	assert.Equal(t, abi.Structs[3].Fields[0].Name, "struct_4_field_1")
	assert.Equal(t, abi.Structs[3].Fields[0].Type, "string")

	// Actions
	assert.Equal(t, abi.Actions[0].Name, chain.N("action_name_1"))
	assert.Equal(t, abi.Actions[0].Type, "struct_name_1")
	assert.Equal(t, abi.Actions[0].RicardianContract, "")

	// Tables
	assert.Equal(t, abi.Tables[0].Name, chain.N("table_name_1"))
	assert.Equal(t, abi.Tables[0].Type, "struct_name_1")
	assert.Equal(t, abi.Tables[0].IndexType, "i64")
	assert.Equal(t, abi.Tables[0].KeyNames[0], "key_name_1")
	assert.Equal(t, abi.Tables[0].KeyNames[1], "key_name_2")
	assert.Equal(t, abi.Tables[0].KeyTypes[0], "string")
	assert.Equal(t, abi.Tables[0].KeyTypes[1], "int")
}

func mockAPI(handler http.HandlerFunc) (*api.Client, *httptest.Server) {
	server := httptest.NewServer(handler)

	return api.New(server.URL), server
}

func TestManager_GetAbiFromCache(t *testing.T) {
	cache := cache.NewCache("thalos::cache::abi::test", cache.NewMemoryStore())

	api, _ := mockAPI(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	mgr := NewAbiManager(cache, api)

	abi := chain.Abi{}
	err := json.Unmarshal([]byte(abiString), &abi)
	assert.NoError(t, err)

	err = mgr.SetAbi(chain.N("testaccount"), &abi)
	assert.NoError(t, err)

	c_abi, err := mgr.GetAbi(chain.N("testaccount"))
	assert.NoError(t, err)
	assert_abi(t, c_abi)
}

func TestManager_GetAbiFromAPI(t *testing.T) {
	cache := cache.NewCache("thalos::cache::abi::test", cache.NewMemoryStore())

	api, _ := mockAPI(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := fmt.Sprintf(`{"account_name": "testaccount", "abi": %s}`, abiString)

		_, err := w.Write([]byte(body))
		assert.NoError(t, err)
	}))

	mgr := NewAbiManager(cache, api)

	c_abi, err := mgr.GetAbi(chain.N("testaccount"))
	assert.NoError(t, err)

	fmt.Println(c_abi)

	assert_abi(t, c_abi)
}

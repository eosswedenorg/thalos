package server

import (
	"testing"

	"github.com/eosswedenorg/thalos/api/message"
	"github.com/eosswedenorg/thalos/internal/filter"
	"github.com/shufflingpixels/antelope-go/chain"
	"github.com/stretchr/testify/require"
)

func TestShipProcessor_ShouldProcessTableDeltaName(t *testing.T) {
	processor := &ShipProcessor{}

	require.True(t, processor.shouldProcessTableDeltaName("contract_row"))
	require.True(t, processor.shouldProcessTableDeltaName("resource_limits"))

	processor.SetTableDeltaWhitelist(*filter.New(map[string][]string{
		"eosio.token": {"accounts"},
	}))

	require.True(t, processor.shouldProcessTableDeltaName("contract_row"))
	require.False(t, processor.shouldProcessTableDeltaName("resource_limits"))
}

func TestShipProcessor_FilterWhitelistedContractRows(t *testing.T) {
	processor := &ShipProcessor{}
	processor.SetTableDeltaWhitelist(*filter.New(map[string][]string{
		"eosio.token": {"accounts"},
		"eosio":       {"*"},
	}))

	rows := []message.TableDeltaRow{
		{Data: map[string]any{"code": "eosio.token", "table": "accounts"}},
		{Data: map[string]any{"code": "eosio.token", "table": "stat"}},
		{Data: map[string]any{"code": "eosio", "table": "users"}},
		{Data: map[string]any{"scope": "missing-code"}},
		{Data: map[string]any{"code": 42, "table": "accounts"}},
		{Data: map[string]any{"code": "eosio.token", "table": 1}},
		{Data: nil},
	}

	actual := processor.filterWhitelistedContractRows(rows)

	require.Len(t, actual, 2)
	require.Equal(t, "eosio.token", actual[0].Data["code"])
	require.Equal(t, "accounts", actual[0].Data["table"])
	require.Equal(t, "eosio", actual[1].Data["code"])
}

func TestShipProcessor_FilterWhitelistedContractRows_StringerCode(t *testing.T) {
	processor := &ShipProcessor{}
	processor.SetTableDeltaWhitelist(*filter.New(map[string][]string{
		"eosio.token": {"accounts"},
	}))

	rows := []message.TableDeltaRow{
		{Data: map[string]any{"code": chain.N("eosio.token"), "table": chain.N("accounts")}},
		{Data: map[string]any{"code": chain.N("eosio.token"), "table": chain.N("stat")}},
	}

	actual := processor.filterWhitelistedContractRows(rows)

	require.Len(t, actual, 1)
	require.Equal(t, chain.N("eosio.token"), actual[0].Data["code"])
}

func TestShipProcessor_FilterWhitelistedContractRows_NoWhitelist(t *testing.T) {
	processor := &ShipProcessor{}

	rows := []message.TableDeltaRow{
		{Data: map[string]any{"code": "eosio.token", "table": "accounts"}},
		{Data: map[string]any{"scope": "missing-code"}},
	}

	actual := processor.filterWhitelistedContractRows(rows)

	require.Equal(t, rows, actual)
}

func TestShipProcessor_SetPostTransactions(t *testing.T) {
	processor := &ShipProcessor{postTransactions: true}

	processor.SetPostTransactions(false)

	require.False(t, processor.postTransactions)
}

func TestShipProcessor_SetPostActions(t *testing.T) {
	processor := &ShipProcessor{postActions: true}

	processor.SetPostActions(false)

	require.False(t, processor.postActions)
}

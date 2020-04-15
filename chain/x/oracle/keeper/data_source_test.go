package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	me "github.com/bandprotocol/bandchain/chain/x/oracle/keeper"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func mockDataSource(ctx sdk.Context, keeper me.Keeper) error {
	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source"
	description := "description"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))

	// Size of "executable" is  10 bytes
	executable := []byte("executable")
	_, err := keeper.AddDataSource(ctx, owner, name, description, fee, executable)
	return err
}

func TestGetterSetterDataSource(t *testing.T) {
	app, ctx := createTestInput()
	keeper := app.OracleKeeper

	_, err := keeper.GetDataSource(ctx, 1)
	require.NotNil(t, err)

	err = mockDataSource(ctx, keeper)
	require.Nil(t, err)

	actualDataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, sdk.AccAddress([]byte("owner")), actualDataSource.Owner)
	require.Equal(t, "data_source", actualDataSource.Name)
	require.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uband", 10)), actualDataSource.Fee)
	require.Equal(t, []byte("executable"), actualDataSource.Executable)
}

func TestAddDataSourceMustReturnCorrectID(t *testing.T) {
	app, ctx := createTestInput()
	keeper := app.OracleKeeper

	_, err := keeper.GetDataSource(ctx, 1)
	require.NotNil(t, err)

	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source"
	description := "description"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	executable := []byte("executable")

	id, err := keeper.AddDataSource(ctx, owner, name, description, fee, executable)
	require.Nil(t, err)
	require.Equal(t, types.DataSourceID(1), id)

	id, err = keeper.AddDataSource(ctx, owner, name, description, fee, executable)
	require.Nil(t, err)
	require.Equal(t, types.DataSourceID(2), id)
}

func TestAddTooLongDataSource(t *testing.T) {
	app, ctx := createTestInput()
	keeper := app.OracleKeeper

	_, err := keeper.GetDataSource(ctx, 1)
	require.NotNil(t, err)

	// Set MaxDataSourceExecutableSize to 20
	keeper.SetParam(ctx, types.KeyMaxDataSourceExecutableSize, 20)

	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source"
	description := "description"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	tooLongExecutable := []byte("The number of bytes of this data source is 80 which is obviously longer than 20.")

	_, err = keeper.AddDataSource(ctx, owner, name, description, fee, tooLongExecutable)
	require.NotNil(t, err)
}

func TestAddTooLongDataSourceName(t *testing.T) {
	app, ctx := createTestInput()
	keeper := app.OracleKeeper

	_, err := keeper.GetDataSource(ctx, 1)
	require.NotNil(t, err)

	// Set MaxNameLength to 5
	keeper.SetParam(ctx, types.KeyMaxNameLength, 5)

	owner := sdk.AccAddress([]byte("owner"))
	tooLongName := "data_source"
	description := "description"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	executable := []byte("executable")

	_, err = keeper.AddDataSource(ctx, owner, tooLongName, description, fee, executable)
	require.NotNil(t, err)
}

func TestAddTooLongDataSourceDescription(t *testing.T) {
	app, ctx := createTestInput()
	keeper := app.OracleKeeper

	_, err := keeper.GetDataSource(ctx, 1)
	require.NotNil(t, err)

	// Set MaxDescriptionLength to 5
	keeper.SetParam(ctx, types.KeyMaxDescriptionLength, 5)

	owner := sdk.AccAddress([]byte("owner"))
	name := "data_source"
	tooLongDescription := "description"
	fee := sdk.NewCoins(sdk.NewInt64Coin("uband", 10))
	executable := []byte("executable")

	_, err = keeper.AddDataSource(ctx, owner, name, tooLongDescription, fee, executable)
	require.NotNil(t, err)
}
func TestEditDataSource(t *testing.T) {
	app, ctx := createTestInput()
	keeper := app.OracleKeeper

	err := mockDataSource(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "data_source_2"
	newDescription := "description_2"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 1))
	newExecutable := []byte("executable_2")

	err = keeper.EditDataSource(ctx, 1, newOwner, newName, newDescription, newFee, newExecutable)
	require.Nil(t, err)

	actualDataSource, err := keeper.GetDataSource(ctx, 1)
	require.Nil(t, err)
	require.Equal(t, newOwner, actualDataSource.Owner)
	require.Equal(t, newName, actualDataSource.Name)
	require.Equal(t, newDescription, actualDataSource.Description)
	require.Equal(t, newFee, actualDataSource.Fee)
	require.Equal(t, newExecutable, actualDataSource.Executable)
}

func TestEditTooLongDataSource(t *testing.T) {
	app, ctx := createTestInput()
	keeper := app.OracleKeeper

	// Set MaxDataSourceExecutableSize to 20
	keeper.SetParam(ctx, types.KeyMaxDataSourceExecutableSize, 20)
	err := mockDataSource(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "data_source_2"
	newDescription := "new_description"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 1))
	newTooLongExecutable := []byte("The number of bytes of this data source is 80 which is obviously longer than 20.")

	err = keeper.EditDataSource(ctx, 1, newOwner, newName, newDescription, newFee, newTooLongExecutable)
	require.NotNil(t, err)
}

func TestEditTooLongDataSourceName(t *testing.T) {
	app, ctx := createTestInput()
	keeper := app.OracleKeeper

	//SetMaxNameLength to 20
	keeper.SetParam(ctx, types.KeyMaxNameLength, 20)
	err := mockDataSource(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newTooLongName := "Toooooooooooo Longggggggggggg"
	newDescription := "new_description"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 1))
	newExecutable := []byte("executable")

	err = keeper.EditDataSource(ctx, 1, newOwner, newTooLongName, newDescription, newFee, newExecutable)
	require.NotNil(t, err)
}

func TestEditTooLongDataSourceDescription(t *testing.T) {
	app, ctx := createTestInput()
	keeper := app.OracleKeeper

	//Set MaxDescriptionLength to 20
	keeper.SetParam(ctx, types.KeyMaxDescriptionLength, 20)
	err := mockDataSource(ctx, keeper)
	require.Nil(t, err)

	newOwner := sdk.AccAddress([]byte("owner2"))
	newName := "data_source"
	newTooLongDescription := "Tooooooooo Loooooooooooooong description"
	newFee := sdk.NewCoins(sdk.NewInt64Coin("uband", 1))
	newExecutable := []byte("executable")

	err = keeper.EditDataSource(ctx, 1, newOwner, newName, newTooLongDescription, newFee, newExecutable)
	require.NotNil(t, err)
}

func TestGetAllDataSources(t *testing.T) {
	app, ctx := createTestInput()
	keeper := app.OracleKeeper

	dataSources := []types.DataSource{
		types.NewDataSource(
			sdk.AccAddress([]byte("owner1")),
			"name1",
			"description1",
			sdk.NewCoins(sdk.NewInt64Coin("uband", 10)),
			[]byte("code1"),
		),
		types.NewDataSource(
			sdk.AccAddress([]byte("owner2")),
			"name2",
			"description2",
			sdk.NewCoins(sdk.NewInt64Coin("uband", 100)),
			[]byte("code2"),
		),
	}
	keeper.SetDataSource(ctx, 1, dataSources[0])
	keeper.SetDataSource(ctx, 2, dataSources[1])

	require.Equal(t, dataSources, keeper.GetAllDataSources(ctx))
}

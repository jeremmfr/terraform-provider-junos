package providerfwk

import (
	"context"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type dataSourceDataFillID interface {
	fillID()
}

type dataSourceDataReadWithoutArg interface {
	dataSourceDataFillID
	read(context.Context, *junos.Session) error
}

type dataSourceDataReadWith1String interface {
	dataSourceDataFillID
	read(context.Context, string, *junos.Session) error
}

type dataSourceDataReadWith1String2Bool interface {
	dataSourceDataFillID
	read(context.Context, string, bool, bool, *junos.Session) error
}

type dataSourceDataFromResource interface {
	copyFromResourceData(any)
}

type junosDataSource interface {
	junosClient() *junos.Client
}

func defaultDataSourceRead(
	ctx context.Context,
	dsc junosDataSource,
	mainAttrValues []any,
	data dataSourceDataFillID,
	resp *datasource.ReadResponse,
) {
	junSess, err := dsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	if data0, ok := data.(dataSourceDataReadWithoutArg); ok {
		err = data0.read(ctx, junSess)
	}
	if data1, ok := data.(dataSourceDataReadWith1String); ok {
		err = data1.read(ctx, mainAttrValues[0].(string), junSess)
	}
	if data1and2, ok := data.(dataSourceDataReadWith1String2Bool); ok {
		err = data1and2.read(ctx, mainAttrValues[0].(string), mainAttrValues[1].(bool), mainAttrValues[2].(bool), junSess)
	}
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ReadErrSummary, err.Error())

		return
	}

	data.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func defaultDataSourceReadFromResource(
	ctx context.Context,
	dsc junosDataSource,
	mainAttrValues []string,
	data dataSourceDataFromResource,
	rscData resourceDataNullID,
	resp *datasource.ReadResponse,
	notFoundDetailMsg string,
) {
	junSess, err := dsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	if data1, ok := rscData.(resourceDataReadFrom1String); ok {
		err = data1.read(ctx, mainAttrValues[0], junSess)
	}
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ReadErrSummary, err.Error())

		return
	}
	if rscData.nullID() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			notFoundDetailMsg,
		)

		return
	}

	data.copyFromResourceData(rscData)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

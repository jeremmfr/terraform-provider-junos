package providerfwk

import (
	"context"
	"fmt"
	"strings"

	"github.com/jeremmfr/terraform-provider-junos/internal/junos"
	"github.com/jeremmfr/terraform-provider-junos/internal/tfdiag"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type resourceDataNullID interface {
	nullID() bool
}

type resourceDataReadWithoutArg interface {
	resourceDataNullID
	read(context.Context, *junos.Session) error
}

type resourceDataReadFrom1String interface {
	resourceDataNullID
	read(context.Context, string, *junos.Session) error
}

type resourceDataReadFrom2String interface {
	resourceDataNullID
	read(context.Context, string, string, *junos.Session) error
}

type resourceDataReadFrom3String interface {
	resourceDataNullID
	read(context.Context, string, string, string, *junos.Session) error
}

type resourceDataReadFrom1String1Bool1String interface {
	resourceDataNullID
	read(context.Context, string, bool, string, *junos.Session) error
}

type resourceDataReadFrom4String interface {
	resourceDataNullID
	read(context.Context, string, string, string, string, *junos.Session) error
}

type resourceDataReadFrom2String1Bool1String interface {
	resourceDataNullID
	read(context.Context, string, string, bool, string, *junos.Session) error
}

// resourceCreateCheck: func to pre and post check when creating a resource
// need to return true if OK and false if NOT OK.
type resourceCreateCheck func(context.Context, *junos.Session) bool

type resourceDataSet interface {
	set(context.Context, *junos.Session) (path.Path, error)
}

type resourceDataFirstSet interface {
	resourceDataSet
	fillID()
}

type resourceDataReadPrivateToState interface {
	readPrivateToState(context.Context, *junos.Session, privateStateSetter) error
}

type privateStateSetter interface {
	SetKey(context.Context, string, []byte) diag.Diagnostics
}

type privateStateGetter interface {
	GetKey(context.Context, string) ([]byte, diag.Diagnostics)
}

type resourceDataDel interface {
	del(context.Context, *junos.Session) error
}

type resourceDataDelWithOpts interface {
	resourceDataDel
	delOpts(context.Context, *junos.Session) error
}

type junosResource interface {
	junosClient() *junos.Client
	typeName() string
}

func defaultResourceCreate(
	ctx context.Context,
	rsc junosResource,
	preCheck resourceCreateCheck,
	postCheck resourceCreateCheck,
	plan resourceDataFirstSet,
	resp *resource.CreateResponse,
) {
	if rsc.junosClient().FakeCreateSetFile() {
		junSess := rsc.junosClient().NewSessionWithoutNetconf(ctx)

		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
			} else {
				resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
			}

			return
		}

		plan.fillID()
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigClearUnlockWarnSummary, junSess.ConfigClear())...)
	}()

	if preCheck != nil && !preCheck(ctx, junSess) {
		return
	}

	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf(ctx, "create resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	if postCheck != nil && !postCheck(ctx, junSess) {
		return
	}

	plan.fillID()
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if dataPriv, ok := plan.(resourceDataReadPrivateToState); ok {
		if err := dataPriv.readPrivateToState(ctx, junSess, resp.Private); err != nil {
			resp.Diagnostics.AddError(tfdiag.ReadPrivateToStateErrSummary, err.Error())
		}
	}
}

func defaultResourceRead(
	ctx context.Context,
	rsc junosResource,
	mainAttrValues []any,
	data resourceDataNullID,
	beforeSetState func(),
	resp *resource.ReadResponse,
) {
	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	junos.MutexLock()
	if data0, ok := data.(resourceDataReadWithoutArg); ok {
		err = data0.read(ctx, junSess)
	}
	if data1, ok := data.(resourceDataReadFrom1String); ok {
		err = data1.read(
			ctx,
			mainAttrValues[0].(string),
			junSess,
		)
	}
	if data2, ok := data.(resourceDataReadFrom2String); ok {
		err = data2.read(
			ctx,
			mainAttrValues[0].(string),
			mainAttrValues[1].(string),
			junSess,
		)
	}
	if data3, ok := data.(resourceDataReadFrom3String); ok {
		err = data3.read(
			ctx,
			mainAttrValues[0].(string),
			mainAttrValues[1].(string),
			mainAttrValues[2].(string),
			junSess,
		)
	}
	if data3, ok := data.(resourceDataReadFrom1String1Bool1String); ok {
		err = data3.read(
			ctx,
			mainAttrValues[0].(string),
			mainAttrValues[1].(bool),
			mainAttrValues[2].(string),
			junSess,
		)
	}
	if data4, ok := data.(resourceDataReadFrom4String); ok {
		err = data4.read(
			ctx,
			mainAttrValues[0].(string),
			mainAttrValues[1].(string),
			mainAttrValues[2].(string),
			mainAttrValues[3].(string),
			junSess,
		)
	}
	if data4, ok := data.(resourceDataReadFrom2String1Bool1String); ok {
		err = data4.read(
			ctx,
			mainAttrValues[0].(string),
			mainAttrValues[1].(string),
			mainAttrValues[2].(bool),
			mainAttrValues[3].(string),
			junSess,
		)
	}
	junos.MutexUnlock()
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}

	if data.nullID() {
		resp.State.RemoveResource(ctx)

		return
	}

	if beforeSetState != nil {
		beforeSetState()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func defaultResourceUpdate(
	ctx context.Context,
	rsc junosResource,
	state resourceDataDel,
	plan resourceDataSet,
	resp *resource.UpdateResponse,
) {
	if rsc.junosClient().FakeUpdateAlso() {
		junSess := rsc.junosClient().NewSessionWithoutNetconf(ctx)

		if stateOpts, ok := state.(resourceDataDelWithOpts); ok {
			if err := stateOpts.delOpts(ctx, junSess); err != nil {
				resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

				return
			}
		} else {
			if err := state.del(ctx, junSess); err != nil {
				resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

				return
			}
		}
		if errPath, err := plan.set(ctx, junSess); err != nil {
			if !errPath.Equal(path.Empty()) {
				resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
			} else {
				resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
			}

			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)

		return
	}

	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigClearUnlockWarnSummary, junSess.ConfigClear())...)
	}()

	if stateOpts, ok := state.(resourceDataDelWithOpts); ok {
		if err := stateOpts.delOpts(ctx, junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

			return
		}
	} else {
		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

			return
		}
	}
	if errPath, err := plan.set(ctx, junSess); err != nil {
		if !errPath.Equal(path.Empty()) {
			resp.Diagnostics.AddAttributeError(errPath, tfdiag.ConfigSetErrSummary, err.Error())
		} else {
			resp.Diagnostics.AddError(tfdiag.ConfigSetErrSummary, err.Error())
		}

		return
	}
	warns, err := junSess.CommitConf(ctx, "update resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if dataPriv, ok := plan.(resourceDataReadPrivateToState); ok {
		if err := dataPriv.readPrivateToState(ctx, junSess, resp.Private); err != nil {
			resp.Diagnostics.AddError(tfdiag.ReadPrivateToStateErrSummary, err.Error())
		}
	}
}

func defaultResourceDelete(
	ctx context.Context,
	rsc junosResource,
	state resourceDataDel,
	resp *resource.DeleteResponse,
) {
	if rsc.junosClient().FakeDeleteAlso() {
		junSess := rsc.junosClient().NewSessionWithoutNetconf(ctx)

		if err := state.del(ctx, junSess); err != nil {
			resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

			return
		}

		return
	}

	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()
	if err := junSess.ConfigLock(ctx); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigLockErrSummary, err.Error())

		return
	}
	defer func() {
		resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigClearUnlockWarnSummary, junSess.ConfigClear())...)
	}()

	if err := state.del(ctx, junSess); err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigDelErrSummary, err.Error())

		return
	}
	warns, err := junSess.CommitConf(ctx, "delete resource "+rsc.typeName())
	resp.Diagnostics.Append(tfdiag.Warns(tfdiag.ConfigCommitWarnSummary, warns)...)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigCommitErrSummary, err.Error())

		return
	}
}

func defaultResourceImportState(
	ctx context.Context,
	rsc junosResource,
	data resourceDataNullID,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
	notFoundDetailMsg string,
) {
	junSess, err := rsc.junosClient().StartNewSession(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.StartSessErrSummary, err.Error())

		return
	}
	defer junSess.Close()

	if data0, ok := data.(resourceDataReadWithoutArg); ok {
		err = data0.read(ctx, junSess)
	}
	if data1, ok := data.(resourceDataReadFrom1String); ok {
		err = data1.read(ctx, req.ID, junSess)
	}
	if data2, ok := data.(resourceDataReadFrom2String); ok {
		idList := strings.Split(req.ID, junos.IDSeparator)
		if len(idList) < 2 {
			resp.Diagnostics.AddError(
				"Bad ID Format",
				fmt.Sprintf("missing element(s) in id with separator %q", junos.IDSeparator),
			)

			return
		}

		err = data2.read(ctx, idList[0], idList[1], junSess)
	}
	if data3, ok := data.(resourceDataReadFrom3String); ok {
		idList := strings.Split(req.ID, junos.IDSeparator)
		if len(idList) < 3 {
			resp.Diagnostics.AddError(
				"Bad ID Format",
				fmt.Sprintf("missing element(s) in id with separator %q", junos.IDSeparator),
			)

			return
		}

		err = data3.read(ctx, idList[0], idList[1], idList[2], junSess)
	}
	if data4, ok := data.(resourceDataReadFrom4String); ok {
		idList := strings.Split(req.ID, junos.IDSeparator)
		if len(idList) < 4 {
			resp.Diagnostics.AddError(
				"Bad ID Format",
				fmt.Sprintf("missing element(s) in id with separator %q", junos.IDSeparator),
			)

			return
		}

		err = data4.read(ctx, idList[0], idList[1], idList[2], idList[3], junSess)
	}
	if err != nil {
		resp.Diagnostics.AddError(tfdiag.ConfigReadErrSummary, err.Error())

		return
	}

	if data.nullID() {
		resp.Diagnostics.AddError(
			tfdiag.NotFoundErrSummary,
			notFoundDetailMsg,
		)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

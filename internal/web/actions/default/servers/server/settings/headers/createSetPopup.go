package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type CreateSetPopupAction struct {
	actionutils.ParentAction
}

func (this *CreateSetPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateSetPopupAction) RunGet(params struct {
	HeaderPolicyId int64
	Type           string
}) {
	this.Data["headerPolicyId"] = params.HeaderPolicyId
	this.Data["type"] = params.Type

	this.Show()
}

func (this *CreateSetPopupAction) RunPost(params struct {
	HeaderPolicyId int64
	Name           string
	Value          string

	StatusListJSON    []byte
	MethodsJSON       []byte
	DomainsJSON       []byte
	ShouldAppend      bool
	DisableRedirect   bool
	ShouldReplace     bool
	ReplaceValuesJSON []byte

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "设置请求Header，HeaderPolicyId:%d, Name:%s, Value:%s", params.HeaderPolicyId, params.Name, params.Value)

	params.Must.
		Field("name", params.Name).
		Require("请输入Header名称")

	configResp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{HeaderPolicyId: params.HeaderPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyConfig := &shared.HTTPHeaderPolicy{}
	err = json.Unmarshal(configResp.HeaderPolicyJSON, policyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// status list
	var statusList = []int32{}
	if len(params.StatusListJSON) > 0 {
		err = json.Unmarshal(params.StatusListJSON, &statusList)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// methods
	var methods = []string{}
	if len(params.MethodsJSON) > 0 {
		err = json.Unmarshal(params.MethodsJSON, &methods)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// domains
	var domains = []string{}
	if len(params.DomainsJSON) > 0 {
		err = json.Unmarshal(params.DomainsJSON, &domains)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// replace values
	var replaceValues = []*shared.HTTPHeaderReplaceValue{}
	if len(params.ReplaceValuesJSON) > 0 {
		err = json.Unmarshal(params.ReplaceValuesJSON, &replaceValues)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// 创建Header
	createHeaderResp, err := this.RPC().HTTPHeaderRPC().CreateHTTPHeader(this.AdminContext(), &pb.CreateHTTPHeaderRequest{
		Name:              params.Name,
		Value:             params.Value,
		Status:            statusList,
		Methods:           methods,
		Domains:           domains,
		ShouldAppend:      params.ShouldAppend,
		DisableRedirect:   params.DisableRedirect,
		ShouldReplace:     params.ShouldReplace,
		ReplaceValuesJSON: params.ReplaceValuesJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	headerId := createHeaderResp.HeaderId

	// 保存
	refs := policyConfig.SetHeaderRefs
	refs = append(refs, &shared.HTTPHeaderRef{
		IsOn:     true,
		HeaderId: headerId,
	})
	refsJSON, err := json.Marshal(refs)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicySettingHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicySettingHeadersRequest{
		HeaderPolicyId: params.HeaderPolicyId,
		HeadersJSON:    refsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}

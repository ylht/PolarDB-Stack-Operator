/* 
*Copyright (c) 2019-2021, Alibaba Group Holding Limited;
*Licensed under the Apache License, Version 2.0 (the "License");
*you may not use this file except in compliance with the License.
*You may obtain a copy of the License at

*   http://www.apache.org/licenses/LICENSE-2.0

*Unless required by applicable law or agreed to in writing, software
*distributed under the License is distributed on an "AS IS" BASIS,
*WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*See the License for the specific language governing permissions and
*limitations under the License.
 */


package workflow_shared

import (
	"context"

	"github.com/go-logr/logr"
	"gitlab.alibaba-inc.com/polar-as/polar-mpd-controller/pkg/define"
	wf "gitlab.alibaba-inc.com/polar-as/polar-mpd-controller/pkg/wfimpl"
	"gitlab.alibaba-inc.com/polar-as/polar-wf-engine/statemachine"
)

func checkSwitchRw(obj statemachine.StateResource) (*statemachine.Event, error) {
	cluster := obj.(*wf.MpdClusterResource).GetMpdCluster()

	if cluster.Annotations != nil {
		insId, ok := cluster.Annotations[define.AnnotationSwitchRw]
		if ok && insId != cluster.Status.LeaderInstanceId {
			return statemachine.CreateEvent(statemachine.EventName(define.WorkflowStateSwitchRw), map[string]interface{}{
				"insId": insId,
			}), nil
		}
	}
	return nil, nil
}

func switchRwMainEnter(obj statemachine.StateResource) error {
	resourceWf, err := wf.GetSharedStorageClusterWfManager().CreateResourceWorkflow(obj)
	if err != nil {
		return err
	}
	err = resourceWf.CommonWorkFlowMainEnter(context.TODO(), obj, "SharedStorageClusterSwitchRw", false, checkSwitchRw)
	return err
}

type SwitchRw struct {
	wf.SharedStorageClusterStepBase
}

func (step *SwitchRw) DoStep(ctx context.Context, logger logr.Logger) error {
	insId, ok := step.Resource.Annotations[define.AnnotationSwitchRw]

	if ok && insId != step.Resource.Status.LeaderInstanceId {
		return step.Service.Switchover(ctx, step.Model, insId)
	}
	return nil
}
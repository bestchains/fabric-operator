# 功能演示

## 准备工作

创立集群，部署 operator，以下内容在代码根目录操作。

<detail>
一些操作可能包括

```bash
# 创建kind集群
$ make kind
# 部署CRD
$ kubectl kustomize config/crd | kubectl apply -f -
```

</detail>

## 组织和联盟管理

### 1. 创建 3 个组织

```bash
kubectl kustomize config/samples/stack/org1 | kubectl apply -f -
kubectl kustomize config/samples/stack/org2 | kubectl apply -f -
kubectl kustomize config/samples/stack/org3 | kubectl apply -f -
```

每个组织在 kubernetes 中对应一个 ns。

<details>

<summary>详细yaml为:</summary>

```bash
$ kubectl get ns
NAME                 STATUS   AGE
default              Active   16m
ingress-nginx        Active   16m
kube-node-lease      Active   16m
kube-public          Active   16m
kube-system          Active   16m
local-path-storage   Active   16m
org1                 Active   16m
org2                 Active   16m
org3                 Active   16m
```

</details>

同时对应一个 Namespaced 纬度的 CRD：Organization ：

<details>

<summary>status 中会显示当前状态，如下:</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Organization
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Organization","metadata":{"annotations":{},"name":"org1","namespace":"org1"},"spec":{"admin":"orgadmin","adminSecret":"org-admin-secret","caReference":{"ca":"ca","name":"org-ca"},"description":"test org","displayName":"test organization","license":{"accept":true}}}
  creationTimestamp: "2022-12-08T04:12:00Z"
  generation: 1
  name: org1
  namespace: org1
  resourceVersion: "1555"
  uid: 6848c469-1b42-4161-8138-f2c803d5fc73
spec:
  admin: orgadmin
  adminSecret: org-admin-secret
  caReference:
    ca: ca
    name: org-ca
  description: test org
  displayName: test organization
  license:
    accept: true
status:
  errorcode: 28
  lastHeartbeatTime: 2022-12-08 12:13:00.644114 +0800 CST m=+1.119948975
  message: 'Code: 28 - failed to initialize organization: failed to get ca connection
    profile: ConfigMap "org-ca-connection-profile" not found'
  reason: errorOccurredDuringReconcile
  status: "True"
  type: Error
```

</details>

### 2. 创建联盟

某个组织要想创建一个联盟 (Federation)，需要按以下步骤：

1. 在自己的 ns 中创建一个联盟 (Federation) 资源。
2. 创建一个提案（Proposal），类型为创建联盟（createFederation），这会触发在每个可以投票的组织中自动创建 选票（vote）资源。
3. 等待每个可以投票的组织手动修改选票（Vote），即投票，表明自己的意愿，同意或者拒绝。
4. 投票符合提案（Proposal）中规定的策略（spec.Policy）后，成功创建联盟。

#### 2.1 创建联盟 federation 的 CRD

```bash
kubectl apply -f config/samples/ibp.com_v1beta1_federation.yaml
```

<details>

<summary>详细yaml为:</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Federation
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Federation","metadata":{"annotations":{},"name":"federation-sample","namespace":"org1"},"spec":{"description":"federation with org0 \u0026 org1","license":{"accept":true},"members":[{"initiator":true,"name":"org1","namespace":"org1"},{"name":"org2","namespace":"org2"}],"policy":"ALL"}}
  creationTimestamp: "2022-12-08T04:14:40Z"
  generation: 1
  name: federation-sample
  namespace: org1
  resourceVersion: "1908"
  uid: 534fb03b-e11f-4f5c-be48-b57c4d71cf04
spec:
  description: federation with org0 & org1
  license:
    accept: true
  members:
  - initiator: true
    name: org1
    namespace: org1
  - name: org2
    namespace: org2
  policy: ALL
status:
  lastHeartbeatTime: 2022-12-08 12:14:41.010557 +0800 CST m=+101.483887881
  status: "True"
  type: FederationPending
```

</details>

#### 2.2 发起提案 proposal，内容为创建联盟

```bash
kubectl apply -f config/samples/ibp.com_v1beta1_proposal_create_federation.yaml
```

<details>

<summary>一个正在投票中的proposal如下：</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Proposal
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Proposal","metadata":{"annotations":{},"name":"create-federation-sample"},"spec":{"createFederation":{},"federation":{"name":"federation-sample","namespace":"org1"},"initiator":{"name":"org1","namespace":"org1"},"policy":"All"}}
  creationTimestamp: "2022-12-08T04:15:20Z"
  generation: 1
  name: create-federation-sample
  resourceVersion: "1972"
  uid: dced6beb-5632-459b-8d2d-6ef67a66b71b
spec:
  createFederation: {}
  deprecated: false
  federation:
    name: federation-sample
    namespace: org1
  initiator:
    name: org1
    namespace: org1
  policy: All
status:
  phase: Voting
  votes:
  - description: ""
    organization:
      name: org2
      namespace: org2
  - decision: true
    description: ""
    organization:
      name: org1
      namespace: org1
    phase: Voted
    startTime: "2022-12-08T04:15:20Z"
```

</details>

#### 2.3 涉及到的组织投票

投票 Vote 是一个 Namespaced 纬度的 CRD。

controller 会在每个有投票权的 Organization 下创建对应的 Vote，本示例中，会在 org1 和 org2 的 ns 下创建 Vote，因为 org1 是本提案的发起人，因此会自动将 org1 中的 Vote 置为 同意，等待 org2 投票:

<details>

<summary>此时全部的vote如下:</summary>

```yaml
apiVersion: v1
items:
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T04:15:20Z"
    generation: 2
    labels:
      app: create-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org1-create-federation-sample
    namespace: org1
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: create-federation-sample
      uid: dced6beb-5632-459b-8d2d-6ef67a66b71b
    resourceVersion: "1971"
    uid: 08e1a28a-d83a-4f2e-9dd1-6b79434c7719
  spec:
    decision: true
    description: ""
    organizationName: org1
    proposalName: create-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T04:15:20Z"
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T04:15:20Z"
    generation: 1
    labels:
      app: create-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org2-create-federation-sample
    namespace: org2
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: create-federation-sample
      uid: dced6beb-5632-459b-8d2d-6ef67a66b71b
    resourceVersion: "1966"
    uid: 45f1a332-1e9e-4ce3-a95a-139e3a302ef8
  spec:
    description: ""
    organizationName: org2
    proposalName: create-federation-sample
  status:
    phase: Created
kind: List
metadata:
  resourceVersion: ""
```

</details>

org2 投同意票可以用以下命令：

```bash
kubectl patch vote -n org2 vote-org2-create-federation-sample --type='json' -p='[{"op": "replace", "path": "/spec/decision", "value": true}]'
```

<details>

<summary>此时vote更新为：</summary>

```yaml
apiVersion: v1
items:
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T04:15:20Z"
    generation: 2
    labels:
      app: create-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org1-create-federation-sample
    namespace: org1
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: create-federation-sample
      uid: dced6beb-5632-459b-8d2d-6ef67a66b71b
    resourceVersion: "1971"
    uid: 08e1a28a-d83a-4f2e-9dd1-6b79434c7719
  spec:
    decision: true
    description: ""
    organizationName: org1
    proposalName: create-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T04:15:20Z"
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T04:15:20Z"
    generation: 2
    labels:
      app: create-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org2-create-federation-sample
    namespace: org2
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: create-federation-sample
      uid: dced6beb-5632-459b-8d2d-6ef67a66b71b
    resourceVersion: "2194"
    uid: 45f1a332-1e9e-4ce3-a95a-139e3a302ef8
  spec:
    decision: true
    description: ""
    organizationName: org2
    proposalName: create-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T04:18:25Z"
kind: List
metadata:
  resourceVersion: ""
```

</details>

#### 2.4 投票成功，联盟创立

因为 所有组织都同意了创建联盟的提案，proposal 的状态更新为 `Finished`:

```yaml
status:
  conditions:
  - lastTransitionTime: "2022-12-08T04:18:25Z"
    message: Success
    reason: Success
    status: "True"
    type: Succeeded
  phase: Finished
  votes:
  - decision: true
    description: ""
    organization:
      name: org2
      namespace: org2
    phase: Voted
    startTime: "2022-12-08T04:18:25Z"
  - decision: true
    description: ""
    organization:
      name: org1
      namespace: org1
    phase: Voted
    startTime: "2022-12-08T04:15:20Z"
```

<details>

<summary>完整proposal yaml为：</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Proposal
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Proposal","metadata":{"annotations":{},"name":"create-federation-sample"},"spec":{"createFederation":{},"federation":{"name":"federation-sample","namespace":"org1"},"initiator":{"name":"org1","namespace":"org1"},"policy":"All"}}
  creationTimestamp: "2022-12-08T04:15:20Z"
  generation: 1
  name: create-federation-sample
  resourceVersion: "2196"
  uid: dced6beb-5632-459b-8d2d-6ef67a66b71b
spec:
  createFederation: {}
  deprecated: false
  federation:
    name: federation-sample
    namespace: org1
  initiator:
    name: org1
    namespace: org1
  policy: All
status:
  conditions:
  - lastTransitionTime: "2022-12-08T04:18:25Z"
    message: Success
    reason: Success
    status: "True"
    type: Succeeded
  phase: Finished
  votes:
  - decision: true
    description: ""
    organization:
      name: org2
      namespace: org2
    phase: Voted
    startTime: "2022-12-08T04:18:25Z"
  - decision: true
    description: ""
    organization:
      name: org1
      namespace: org1
    phase: Voted
    startTime: "2022-12-08T04:15:20Z"
```

</details>

同时，联盟 federation 的状态也会同步更新为 `FederationActivated`：

```yaml
status:
  lastHeartbeatTime: 2022-12-08 12:35:44.826345 +0800 CST m=+33.724520930
  status: "True"
  type: FederationActivated
```

<details>

<summary>完整proposal yaml为：</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Federation
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Federation","metadata":{"annotations":{},"name":"federation-sample","namespace":"org1"},"spec":{"description":"federation with org0 \u0026 org1","license":{"accept":true},"members":[{"initiator":true,"name":"org1","namespace":"org1"},{"name":"org2","namespace":"org2"}],"policy":"ALL"}}
  creationTimestamp: "2022-12-08T04:14:40Z"
  generation: 1
  name: federation-sample
  namespace: org1
  resourceVersion: "3510"
  uid: 534fb03b-e11f-4f5c-be48-b57c4d71cf04
spec:
  description: federation with org0 & org1
  license:
    accept: true
  members:
    - initiator: true
      name: org1
      namespace: org1
    - name: org2
      namespace: org2
  policy: ALL
status:
  lastHeartbeatTime: 2022-12-08 12:35:44.826345 +0800 CST m=+33.724520930
  status: "True"
  type: FederationActivated
```

</details>

至此，联盟创立成功。

### 3. 向一个联盟中添加组织

某个组织要在现有的联盟中添加一个组织作为成员，需要按以下步骤：

1. 创建一个提案（Proposal），类型为添加成员（AddMember），这会触发在每个可以投票的组织中自动创建 选票（vote）资源。
2. 等待每个可以投票的组织手动修改选票（Vote），即投票，表明自己的意愿，同意或者拒绝。
3. 投票符合提案（Proposal）中规定的策略（spec.Policy）后，成功添加该组织作为成员。

#### 3.1 发起提案 proposal，内容为添加成员

```bash
kubectl apply -f config/samples/ibp.com_v1beta1_proposal_add_member.yaml
```

<details>

<summary>完整proposal yaml为：</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Proposal
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Proposal","metadata":{"annotations":{},"name":"add-member-federation-sample"},"spec":{"addMember":{"members":[{"name":"org3","namespace":"org3"}]},"federation":{"name":"federation-sample","namespace":"org1"},"initiator":{"name":"org1","namespace":"org1"},"policy":"All"}}
  creationTimestamp: "2022-12-08T05:21:33Z"
  generation: 2
  name: add-member-federation-sample
  resourceVersion: "6911"
  uid: 51ccea30-88f4-4860-8008-9f4d3ec306df
spec:
  addMember:
    members:
    - name: org3
      namespace: org3
  deprecated: false
  endAt: "2022-12-09T05:21:33Z"
  federation:
    name: federation-sample
    namespace: org1
  initiator:
    name: org1
    namespace: org1
  policy: All
  startAt: "2022-12-08T05:21:33Z"
status:
  phase: Voting
  votes:
  - description: ""
    organization:
      name: org3
      namespace: org3
  - decision: true
    description: ""
    organization:
      name: org1
      namespace: org1
    phase: Voted
    startTime: "2022-12-08T05:21:34Z"
  - description: ""
    organization:
      name: org2
      namespace: org2
```

</details>

<details>

<summary>此时选票结果为：</summary>

```bash
$kubectl get vote -A
NAMESPACE   NAME                                     AGE
org1        vote-org1-add-member-federation-sample   67s
org1        vote-org1-create-federation-sample       16m
org2        vote-org2-add-member-federation-sample   67s
org2        vote-org2-create-federation-sample       16m
org3        vote-org3-add-member-federation-sample   67s
```

</details>

#### 3.2 涉及到的组织投票

```bash
kubectl patch vote -n org2 vote-org2-add-member-federation-sample --type='json' -p='[{"op": "replace", "path": "/spec/decision", "value": true}]'
kubectl patch vote -n org3 vote-org3-add-member-federation-sample --type='json' -p='[{"op": "replace", "path": "/spec/decision", "value": true}]'
```

<details>
<summary>此时选票结果为：</summary>

```yaml
apiVersion: v1
items:
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:21:33Z"
    generation: 2
    labels:
      app: add-member-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org1-add-member-federation-sample
    namespace: org1
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: add-member-federation-sample
      uid: 51ccea30-88f4-4860-8008-9f4d3ec306df
    resourceVersion: "6910"
    uid: 1cf0499c-235a-4193-abfb-3054a41c62f0
  spec:
    decision: true
    description: ""
    organizationName: org1
    proposalName: add-member-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T05:21:34Z"
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:06:10Z"
    generation: 2
    labels:
      app: create-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org1-create-federation-sample
    namespace: org1
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: create-federation-sample
      uid: 76060189-0827-4143-96a2-6c988b1a100d
    resourceVersion: "5775"
    uid: a6bcb00f-12ef-4380-8fcd-3e085306bc24
  spec:
    decision: true
    description: ""
    organizationName: org1
    proposalName: create-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T05:06:10Z"
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:21:33Z"
    generation: 2
    labels:
      app: add-member-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org2-add-member-federation-sample
    namespace: org2
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: add-member-federation-sample
      uid: 51ccea30-88f4-4860-8008-9f4d3ec306df
    resourceVersion: "7164"
    uid: 437f9bd6-5f6c-4079-878f-d271a1a65e9d
  spec:
    decision: true
    description: ""
    organizationName: org2
    proposalName: add-member-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T05:25:02Z"
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:06:10Z"
    generation: 2
    labels:
      app: create-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org2-create-federation-sample
    namespace: org2
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: create-federation-sample
      uid: 76060189-0827-4143-96a2-6c988b1a100d
    resourceVersion: "5793"
    uid: ba193749-2ba3-4d89-9bfb-6fcd466c7e9c
  spec:
    decision: true
    description: ""
    organizationName: org2
    proposalName: create-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T05:06:21Z"
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:21:33Z"
    generation: 2
    labels:
      app: add-member-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org3-add-member-federation-sample
    namespace: org3
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: add-member-federation-sample
      uid: 51ccea30-88f4-4860-8008-9f4d3ec306df
    resourceVersion: "7167"
    uid: 0eda44c9-57cd-4b3c-ace5-c274366b77d0
  spec:
    decision: true
    description: ""
    organizationName: org3
    proposalName: add-member-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T05:25:03Z"
kind: List
metadata:
  resourceVersion: ""
```

</details>

<details>
<summary>此时提案结果为：</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Proposal
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Proposal","metadata":{"annotations":{},"name":"add-member-federation-sample"},"spec":{"addMember":{"members":[{"name":"org3","namespace":"org3"}]},"federation":{"name":"federation-sample","namespace":"org1"},"initiator":{"name":"org1","namespace":"org1"},"policy":"All"}}
  creationTimestamp: "2022-12-08T05:21:33Z"
  generation: 2
  name: add-member-federation-sample
  resourceVersion: "7169"
  uid: 51ccea30-88f4-4860-8008-9f4d3ec306df
spec:
  addMember:
    members:
    - name: org3
      namespace: org3
  deprecated: false
  endAt: "2022-12-09T05:21:33Z"
  federation:
    name: federation-sample
    namespace: org1
  initiator:
    name: org1
    namespace: org1
  policy: All
  startAt: "2022-12-08T05:21:33Z"
status:
  conditions:
  - lastTransitionTime: "2022-12-08T05:25:03Z"
    message: Success
    reason: Success
    status: "True"
    type: Succeeded
  phase: Finished
  votes:
  - decision: true
    description: ""
    organization:
      name: org3
      namespace: org3
    phase: Voted
    startTime: "2022-12-08T05:25:03Z"
  - decision: true
    description: ""
    organization:
      name: org1
      namespace: org1
    phase: Voted
    startTime: "2022-12-08T05:21:34Z"
  - decision: true
    description: ""
    organization:
      name: org2
      namespace: org2
    phase: Voted
    startTime: "2022-12-08T05:25:02Z"
```

</details>

#### 3.3 投票成功，成员成功添加

<details>

<summary>联盟yaml为：</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Federation
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Federation","metadata":{"annotations":{},"name":"federation-sample","namespace":"org1"},"spec":{"description":"federation with org0 \u0026 org1","license":{"accept":true},"members":[{"initiator":true,"name":"org1","namespace":"org1"},{"name":"org2","namespace":"org2"}],"policy":"ALL"}}
  creationTimestamp: "2022-12-08T04:14:40Z"
  generation: 2
  name: federation-sample
  namespace: org1
  resourceVersion: "8794"
  uid: 534fb03b-e11f-4f5c-be48-b57c4d71cf04
spec:
  description: federation with org0 & org1
  license:
    accept: true
  members:
  - initiator: true
    name: org1
    namespace: org1
  - name: org2
    namespace: org2
  - name: org3
    namespace: org3
  policy: ALL
status:
  lastHeartbeatTime: 2022-12-08 12:35:44.826345 +0800 CST m=+33.724520930
  status: "True"
  type: FederationActivated
```

</details>

### 4. 从一个联盟中驱逐一个组织

某个组织要在现有的联盟中驱逐一个组织，需要按以下步骤：

1. 创建一个提案（Proposal），类型为驱逐成员（deleteMember），这会触发在每个可以投票的组织中自动创建 选票（vote）资源。
2. 等待每个可以投票的组织手动修改选票（Vote），即投票，表明自己的意愿，同意或者拒绝。
3. 投票符合提案（Proposal）中规定的策略（spec.Policy）后，成功添加该组织作为成员。

#### 4.1 发起提案 proposal，内容为驱逐成员

```bash
kubectl apply -f config/samples/ibp.com_v1beta1_proposal_delete_member.yaml
```

<details>

<summary>完整proposal yaml为：</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Proposal
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Proposal","metadata":{"annotations":{},"name":"delete-member-federation-sample"},"spec":{"deleteMember":{"member":{"name":"org2","namespace":"org2"}},"federation":{"name":"federation-sample","namespace":"org1"},"initiator":{"name":"org1","namespace":"org1"},"policy":"All"}}
  creationTimestamp: "2022-12-08T05:49:08Z"
  generation: 2
  name: delete-member-federation-sample
  resourceVersion: "8967"
  uid: 02ef7287-0857-471e-b2d5-6ce4326c35f4
spec:
  deleteMember:
    member:
      name: org2
      namespace: org2
  deprecated: false
  endAt: "2022-12-09T05:49:08Z"
  federation:
    name: federation-sample
    namespace: org1
  initiator:
    name: org1
    namespace: org1
  policy: All
  startAt: "2022-12-08T05:49:08Z"
status:
  phase: Voting
  votes:
  - decision: true
    description: ""
    organization:
      name: org1
      namespace: org1
    phase: Voted
    startTime: "2022-12-08T05:49:08Z"
  - description: ""
    organization:
      name: org3
      namespace: org3
```

</details>

当驱逐 org2 时，org2 不会拥有选票 Vote:

<details>

<summary>此时选票结果为：</summary>

```yaml
apiVersion: v1
items:
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:49:08Z"
    generation: 2
    labels:
      app: delete-member-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org1-delete-member-federation-sample
    namespace: org1
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: delete-member-federation-sample
      uid: 02ef7287-0857-471e-b2d5-6ce4326c35f4
    resourceVersion: "8966"
    uid: 03235fda-d032-4125-8f99-6ea94e92b93c
  spec:
    decision: true
    description: ""
    organizationName: org1
    proposalName: delete-member-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T05:49:08Z"
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:49:08Z"
    generation: 1
    labels:
      app: delete-member-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org3-delete-member-federation-sample
    namespace: org3
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: delete-member-federation-sample
      uid: 02ef7287-0857-471e-b2d5-6ce4326c35f4
    resourceVersion: "8965"
    uid: f1a05522-7ba2-4363-b992-12c29c9dbda1
  spec:
    description: ""
    organizationName: org3
    proposalName: delete-member-federation-sample
  status:
    phase: Created
kind: List
metadata:
  resourceVersion: ""
```

</details>

#### 4.2 涉及到的组织投票

```bash
kubectl patch vote -n org3 vote-org3-delete-member-federation-sample --type='json' -p='[{"op": "replace", "path": "/spec/decision", "value": true}]'
```

<details>

<summary>此时选票结果为：</summary>

```yaml
apiVersion: v1
items:
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:49:08Z"
    generation: 2
    labels:
      app: delete-member-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org1-delete-member-federation-sample
    namespace: org1
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: delete-member-federation-sample
      uid: 02ef7287-0857-471e-b2d5-6ce4326c35f4
    resourceVersion: "8966"
    uid: 03235fda-d032-4125-8f99-6ea94e92b93c
  spec:
    decision: true
    description: ""
    organizationName: org1
    proposalName: delete-member-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T05:49:08Z"
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:49:08Z"
    generation: 2
    labels:
      app: delete-member-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org3-delete-member-federation-sample
    namespace: org3
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: delete-member-federation-sample
      uid: 02ef7287-0857-471e-b2d5-6ce4326c35f4
    resourceVersion: "9170"
    uid: f1a05522-7ba2-4363-b992-12c29c9dbda1
  spec:
    decision: true
    description: ""
    organizationName: org3
    proposalName: delete-member-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T05:51:45Z"
kind: List
metadata:
  resourceVersion: ""
```

</details>

#### 4.3 投票成功，成员成功驱逐

<details>

<summary>联盟yaml为：</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Federation
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Federation","metadata":{"annotations":{},"name":"federation-sample","namespace":"org1"},"spec":{"description":"federation with org0 \u0026 org1","license":{"accept":true},"members":[{"initiator":true,"name":"org1","namespace":"org1"},{"name":"org2","namespace":"org2"}],"policy":"ALL"}}
  creationTimestamp: "2022-12-08T04:14:40Z"
  generation: 2
  name: federation-sample
  namespace: org1
  resourceVersion: "8794"
  uid: 534fb03b-e11f-4f5c-be48-b57c4d71cf04
spec:
  description: federation with org0 & org1
  license:
    accept: true
  members:
  - initiator: true
    name: org1
    namespace: org1 
  - name: org3
    namespace: org3
  policy: ALL
status:
  lastHeartbeatTime: 2022-12-08 12:35:44.826345 +0800 CST m=+33.724520930
  status: "True"
  type: FederationActivated
```

</details>

### 5. 解散联盟

某个组织想要解散现有联盟，需要按以下步骤：

1. 创建一个提案（Proposal），类型为解散联盟（dissolveFederation），这会触发在每个可以投票的组织中自动创建 选票（vote）资源。
2. 等待每个可以投票的组织手动修改选票（Vote），即投票，表明自己的意愿，同意或者拒绝。
3. 投票符合提案（Proposal）中规定的策略（spec.Policy）后，成功添加该组织作为成员。

#### 5.1 发起提案 proposal，内容为解散联盟

```bash
kubectl apply -f config/samples/ibp.com_v1beta1_proposal_dissolve_federation.yaml
```

<details>

<summary>完整proposal yaml为：</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Proposal
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Proposal","metadata":{"annotations":{},"name":"dissolve-federation-sample"},"spec":{"dissolveFederation":{},"federation":{"name":"federation-sample","namespace":"org1"},"initiator":{"name":"org1","namespace":"org1"},"policy":"All"}}
  creationTimestamp: "2022-12-08T05:58:32Z"
  generation: 2
  name: dissolve-federation-sample
  resourceVersion: "9710"
  uid: 9cd71ce9-22d4-4b4f-b3d4-947e57fdd10d
spec:
  deprecated: false
  dissolveFederation: {}
  endAt: "2022-12-09T05:58:32Z"
  federation:
    name: federation-sample
    namespace: org1
  initiator:
    name: org1
    namespace: org1
  policy: All
  startAt: "2022-12-08T05:58:32Z"
status:
  phase: Voting
  votes:
  - description: ""
    organization:
      name: org3
      namespace: org3
  - decision: true
    description: ""
    organization:
      name: org1
      namespace: org1
    phase: Voted
    startTime: "2022-12-08T05:58:32Z"
```

</details>

<details>

<summary>此时选票结果为：</summary>

```yaml
apiVersion: v1
items:
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:58:32Z"
    generation: 2
    labels:
      app: dissolve-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org1-dissolve-federation-sample
    namespace: org1
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: dissolve-federation-sample
      uid: 9cd71ce9-22d4-4b4f-b3d4-947e57fdd10d
    resourceVersion: "9709"
    uid: 96c880a1-71c7-4af3-92b9-aba2a4904f6c
  spec:
    decision: true
    description: ""
    organizationName: org1
    proposalName: dissolve-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T05:58:32Z"
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:58:32Z"
    generation: 1
    labels:
      app: dissolve-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org3-dissolve-federation-sample
    namespace: org3
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: dissolve-federation-sample
      uid: 9cd71ce9-22d4-4b4f-b3d4-947e57fdd10d
    resourceVersion: "9704"
    uid: 5a537d77-167f-41ac-8ac0-c23946b820fa
  spec:
    description: ""
    organizationName: org3
    proposalName: dissolve-federation-sample
  status:
    phase: Created
kind: List
metadata:
  resourceVersion: ""
```

</details>

#### 5.2 涉及到的组织投票

```bash
kubectl patch vote -n org3 vote-org3-dissolve-federation-sample --type='json' -p='[{"op": "replace", "path": "/spec/decision", "value": true}]'
```

<details>

<summary>此时选票结果为：</summary>

```yaml
apiVersion: v1
items:
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:58:32Z"
    generation: 2
    labels:
      app: dissolve-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org1-dissolve-federation-sample
    namespace: org1
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: dissolve-federation-sample
      uid: 9cd71ce9-22d4-4b4f-b3d4-947e57fdd10d
    resourceVersion: "9709"
    uid: 96c880a1-71c7-4af3-92b9-aba2a4904f6c
  spec:
    decision: true
    description: ""
    organizationName: org1
    proposalName: dissolve-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T05:58:32Z"
- apiVersion: ibp.com/v1beta1
  kind: Vote
  metadata:
    creationTimestamp: "2022-12-08T05:58:32Z"
    generation: 2
    labels:
      app: dissolve-federation-sample
      app.kubernetes.io/instance: fabricproposal
      app.kubernetes.io/managed-by: fabric-operator
      app.kubernetes.io/name: fabric
      creator: fabric
      helm.sh/chart: ibm-fabric
      release: operator
    name: vote-org3-dissolve-federation-sample
    namespace: org3
    ownerReferences:
    - apiVersion: ibp.com/v1beta1
      blockOwnerDeletion: true
      controller: true
      kind: Proposal
      name: dissolve-federation-sample
      uid: 9cd71ce9-22d4-4b4f-b3d4-947e57fdd10d
    resourceVersion: "9832"
    uid: 5a537d77-167f-41ac-8ac0-c23946b820fa
  spec:
    decision: true
    description: ""
    organizationName: org3
    proposalName: dissolve-federation-sample
  status:
    phase: Voted
    startTime: "2022-12-08T06:00:04Z"
kind: List
metadata:
  resourceVersion: ""

```

</details>

#### 5.3 投票成功，联盟成功解散

联盟状态变为：

```yaml
status:
  lastHeartbeatTime: 2022-12-08 14:00:04.395943 +0800 CST m=+302.039719659
  status: "True"
  type: FederationDissolved
```

<details>

<summary>联盟yaml为：</summary>

```yaml
apiVersion: ibp.com/v1beta1
kind: Federation
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ibp.com/v1beta1","kind":"Federation","metadata":{"annotations":{},"name":"federation-sample","namespace":"org1"},"spec":{"description":"federation with org0 \u0026 org1","license":{"accept":true},"members":[{"initiator":true,"name":"org1","namespace":"org1"},{"name":"org2","namespace":"org2"}],"policy":"ALL"}}
  creationTimestamp: "2022-12-08T04:14:40Z"
  generation: 4
  name: federation-sample
  namespace: org1
  resourceVersion: "9841"
  uid: 534fb03b-e11f-4f5c-be48-b57c4d71cf04
spec:
  description: federation with org0 & org1
  license:
    accept: true
  members:
  - initiator: true
    name: org1
    namespace: org1
  policy: ALL
status:
  lastHeartbeatTime: 2022-12-08 14:00:04.395943 +0800 CST m=+302.039719659
  status: "True"
  type: FederationDissolved

```

</details>

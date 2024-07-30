<p>Packages:</p>
<ul>
<li>
<a href="#infrastructure.cluster.x-k8s.io%2fv1alpha1">infrastructure.cluster.x-k8s.io/v1alpha1</a>
</li>
</ul>
<h2 id="infrastructure.cluster.x-k8s.io/v1alpha1">infrastructure.cluster.x-k8s.io/v1alpha1</h2>
<div>
<p>Package v1alpha1 contains API Schema definitions for the settings.gardener.cloud API group</p>
</div>
Resource Types:
<ul></ul>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.MetalCluster">MetalCluster
</h3>
<div>
<p>MetalCluster is the Schema for the metalclusters API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalClusterSpec">
MetalClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
sigs.k8s.io/cluster-api/api/v1beta1.APIEndpoint
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalClusterStatus">
MetalClusterStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.MetalClusterSpec">MetalClusterSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalCluster">MetalCluster</a>)
</p>
<div>
<p>MetalClusterSpec defines the desired state of MetalCluster</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
sigs.k8s.io/cluster-api/api/v1beta1.APIEndpoint
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.MetalClusterStatus">MetalClusterStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalCluster">MetalCluster</a>)
</p>
<div>
<p>MetalClusterStatus defines the observed state of MetalCluster</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ready denotes that the cluster (infrastructure) is ready.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
sigs.k8s.io/cluster-api/api/v1beta1.Conditions
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the MetalCluster.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachine">MetalMachine
</h3>
<div>
<p>MetalMachine is the Schema for the metalmachines API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineSpec">
MetalMachineSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>serverSelector</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.ServerSelector">
ServerSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServerSelector specifies matching criteria for labels on Servers.
This is used to claim specific Server types for a MetalMachine.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineStatus">
MetalMachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineSpec">MetalMachineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachine">MetalMachine</a>, <a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineTemplateResource">MetalMachineTemplateResource</a>)
</p>
<div>
<p>MetalMachineSpec defines the desired state of MetalMachine</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>serverSelector</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.ServerSelector">
ServerSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServerSelector specifies matching criteria for labels on Servers.
This is used to claim specific Server types for a MetalMachine.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineStatus">MetalMachineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachine">MetalMachine</a>)
</p>
<div>
<p>MetalMachineStatus defines the observed state of MetalMachine</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ready indicates the Machine infrastructure has been provisioned and is ready.</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
sigs.k8s.io/cluster-api/errors.MachineStatusError
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set in the event that there is a terminal problem
reconciling the Machine and will contain a succinct value suitable
for machine interpretation.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of Machines
can be added as events to the Machine object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the Machine and will contain a more verbose string suitable
for logging and human consumption.</p>
<p>This field should not be set for transitive errors that a controller
faces that are expected to be fixed automatically over
time (like service outages), but instead indicate that something is
fundamentally wrong with the Machine&rsquo;s spec or the configuration of
the controller, and that manual intervention is required. Examples
of terminal errors would be invalid combinations of settings in the
spec, values that are unsupported by the controller, or the
responsible controller itself being critically misconfigured.</p>
<p>Any transient errors that occur during the reconciliation of Machines
can be added as events to the Machine object and/or logged in the
controller&rsquo;s output.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineTemplate">MetalMachineTemplate
</h3>
<div>
<p>MetalMachineTemplate is the Schema for the metalmachinetemplates API</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineTemplateSpec">
MetalMachineTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineTemplateResource">
MetalMachineTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineTemplateResource">MetalMachineTemplateResource
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineTemplateSpec">MetalMachineTemplateSpec</a>)
</p>
<div>
<p>MetalMachineTemplateResource defines the spec and metadata for MetalMachineTemplate supported by capi.</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
sigs.k8s.io/cluster-api/api/v1beta1.ObjectMeta
</em>
</td>
<td>
<em>(Optional)</em>
<p>Standard object&rsquo;s metadata.
More info: <a href="https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata">https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata</a></p>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineSpec">
MetalMachineSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>serverSelector</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.ServerSelector">
ServerSelector
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServerSelector specifies matching criteria for labels on Servers.
This is used to claim specific Server types for a MetalMachine.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineTemplateSpec">MetalMachineTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineTemplate">MetalMachineTemplate</a>)
</p>
<div>
<p>MetalMachineTemplateSpec defines the desired state of MetalMachineTemplate</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineTemplateResource">
MetalMachineTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.ServerSelector">ServerSelector
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.MetalMachineSpec">MetalMachineSpec</a>)
</p>
<div>
<p>ServerSelector specifies matching criteria for labels on Server.
This is used to claim specific Server types for a Machine</p>
</div>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>matchLabels</code><br/>
<em>
map[string]string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Key/value pairs of labels that must exist on a chosen Server</p>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
</em></p>

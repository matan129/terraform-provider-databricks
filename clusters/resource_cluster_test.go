package clusters

import (
	"fmt"
	"strings"
	"testing"

	"github.com/databrickslabs/terraform-provider-databricks/common"
	"github.com/databrickslabs/terraform-provider-databricks/libraries"

	"github.com/databrickslabs/terraform-provider-databricks/qa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceClusterCreate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/create",
				ExpectedRequest: Cluster{
					NumWorkers:             100,
					ClusterName:            "Shared Autoscaling",
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
					AutoterminationMinutes: 15,
				},
				Response: ClusterInfo{
					ClusterID: "abc",
					State:     ClusterStateRunning,
				},
			},
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/clusters/get?cluster_id=abc",
				Response: ClusterInfo{
					ClusterID:              "abc",
					NumWorkers:             100,
					ClusterName:            "Shared Autoscaling",
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
					AutoterminationMinutes: 15,
					State:                  ClusterStateRunning,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/events",
				ExpectedRequest: EventsRequest{
					ClusterID:  "abc",
					Limit:      1,
					Order:      SortDescending,
					EventTypes: []ClusterEventType{EvTypePinned, EvTypeUnpinned},
				},
				Response: EventsResponse{
					Events:     []ClusterEvent{},
					TotalCount: 0,
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{},
				},
			},
		},
		Create:   true,
		Resource: ResourceCluster(),
		State: map[string]interface{}{
			"autotermination_minutes": 15,
			"cluster_name":            "Shared Autoscaling",
			"spark_version":           "7.1-scala12",
			"node_type_id":            "i3.xlarge",
			"num_workers":             100,
			"is_pinned":               false,
		},
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc", d.Id())
}

func TestResourceClusterCreatePinned(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/create",
				ExpectedRequest: Cluster{
					NumWorkers:             100,
					ClusterName:            "Shared Autoscaling",
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
					AutoterminationMinutes: 15,
				},
				Response: ClusterInfo{
					ClusterID: "abc",
					State:     ClusterStateRunning,
				},
			},
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/clusters/get?cluster_id=abc",
				Response: ClusterInfo{
					ClusterID:              "abc",
					NumWorkers:             100,
					ClusterName:            "Shared Autoscaling",
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
					AutoterminationMinutes: 15,
					State:                  ClusterStateRunning,
				},
			},
			{
				Method:          "POST",
				Resource:        "/api/2.0/clusters/pin",
				ExpectedRequest: ClusterID{ClusterID: "abc"},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{},
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/events",
				ExpectedRequest: EventsRequest{
					ClusterID:  "abc",
					Limit:      1,
					Order:      SortDescending,
					EventTypes: []ClusterEventType{EvTypePinned, EvTypeUnpinned},
				},
				Response: EventsResponse{
					Events: []ClusterEvent{
						{
							ClusterID: "abc",
							Timestamp: int64(123),
							Type:      EvTypePinned,
							Details:   EventDetails{},
						},
					},
					TotalCount: 1,
				},
			},
		},
		Create:   true,
		Resource: ResourceCluster(),
		State: map[string]interface{}{
			"autotermination_minutes": 15,
			"cluster_name":            "Shared Autoscaling",
			"spark_version":           "7.1-scala12",
			"node_type_id":            "i3.xlarge",
			"num_workers":             100,
			"is_pinned":               true,
		},
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc", d.Id())
}

func TestResourceClusterCreate_WithLibraries(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/create",
				ExpectedRequest: Cluster{
					NumWorkers:             100,
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
					AutoterminationMinutes: 60,
				},
				Response: ClusterInfo{
					ClusterID: "abc",
					State:     ClusterStateRunning,
				},
			},
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/clusters/get?cluster_id=abc",
				Response: ClusterInfo{
					ClusterID:              "abc",
					NumWorkers:             100,
					ClusterName:            "Shared Autoscaling",
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
					AutoterminationMinutes: 15,
					State:                  ClusterStateRunning,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/events",
				ExpectedRequest: EventsRequest{
					ClusterID:  "abc",
					Limit:      1,
					Order:      SortDescending,
					EventTypes: []ClusterEventType{EvTypePinned, EvTypeUnpinned},
				},
				Response: EventsResponse{
					Events:     []ClusterEvent{},
					TotalCount: 0,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/libraries/install",
				ExpectedRequest: libraries.ClusterLibraryList{
					ClusterID: "abc",
					Libraries: []libraries.Library{
						{
							Jar: "dbfs://foo.jar",
						},
						{
							Cran: &libraries.Cran{
								Package: "rkeops",
								Repo:    "internal",
							},
						},
						{
							Egg: "dbfs://bar.egg",
						},
						{
							Maven: &libraries.Maven{
								Coordinates: "foo:bar:baz:0.1.0",
								Exclusions:  []string{"org.apache:flink:base"},
								Repo:        "s3://maven-repo-in-s3/release",
							},
						},
						{
							Pypi: &libraries.PyPi{
								Package: "seaborn==1.2.4",
							},
						},
						{
							Whl: "dbfs://baz.whl",
						},
					},
				},
			},
			{
				Method: "GET",
				// 1 of 3 requests
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{
						{
							Library: &libraries.Library{
								Pypi: &libraries.PyPi{
									Package: "seaborn==1.2.4",
								},
							},
							Status: "PENDING",
						},
						{
							Library: &libraries.Library{
								Whl: "dbfs://baz.whl",
							},
							Status: "INSTALLED",
						},
					},
				},
			},
			{
				Method: "GET",
				// 2 of 3 requests
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{
						{
							Library: &libraries.Library{
								Pypi: &libraries.PyPi{
									Package: "seaborn==1.2.4",
								},
							},
							Status: "INSTALLED",
						},
						{
							Library: &libraries.Library{
								Whl: "dbfs://baz.whl",
							},
							Status: "INSTALLED",
						},
					},
				},
			},
			{
				Method: "GET",
				// 3 of 3 requests
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{
						{
							Library: &libraries.Library{
								Pypi: &libraries.PyPi{
									Package: "seaborn==1.2.4",
								},
							},
							Status: "INSTALLED",
						},
						{
							Library: &libraries.Library{
								Whl: "dbfs://baz.whl",
							},
							Status: "INSTALLED",
						},
					},
				},
			},
		},
		Create:   true,
		Resource: ResourceCluster(),
		HCL: `num_workers = 100
		spark_version = "7.1-scala12"
		node_type_id = "i3.xlarge"

		library {
			jar = "dbfs://foo.jar"
		}

		library {
			egg = "dbfs://bar.egg"
		}

		library {
			whl = "dbfs://baz.whl"
		}

		library {
			pypi {
				package = "seaborn==1.2.4"
			}
		}

		library {
			maven {
				coordinates = "foo:bar:baz:0.1.0"
				repo = "s3://maven-repo-in-s3/release"
				exclusions = [
					"org.apache:flink:base"
				]
			}
		}

		library {
			cran {
				package = "rkeops"
				repo = "internal"
			}
		}`,
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc", d.Id())
}

func TestResourceClusterCreate_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/create",
				Response: common.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		Create:   true,
		Resource: ResourceCluster(),
		State: map[string]interface{}{
			"autotermination_minutes": 15,
			"cluster_name":            "Shared Autoscaling",
			"spark_version":           "7.1-scala12",
			"node_type_id":            "i3.xlarge",
			"num_workers":             100,
		},
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "", d.Id(), "Id should be empty for error creates")
}

func TestResourceClusterRead(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/clusters/get?cluster_id=abc",
				Response: ClusterInfo{
					ClusterID:              "abc",
					NumWorkers:             100,
					ClusterName:            "Shared Autoscaling",
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
					AutoterminationMinutes: 15,
					State:                  ClusterStateRunning,
					AutoScale: &AutoScale{
						MaxWorkers: 4,
					},
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/events",
				ExpectedRequest: EventsRequest{
					ClusterID:  "abc",
					Limit:      1,
					Order:      SortDescending,
					EventTypes: []ClusterEventType{EvTypePinned, EvTypeUnpinned},
				},
				Response: EventsResponse{
					Events:     []ClusterEvent{},
					TotalCount: 0,
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{
						{
							Library: &libraries.Library{
								Pypi: &libraries.PyPi{
									Package: "requests",
								},
							},
							Status: "INSTALLED",
						},
					},
				},
			},
		},
		Resource: ResourceCluster(),
		Read:     true,
		ID:       "abc",
		New:      true,
	}.Apply(t)
	require.NoError(t, err, err)
	assert.Equal(t, "abc", d.Id(), "Id should not be empty")
	assert.Equal(t, 15, d.Get("autotermination_minutes"))
	assert.Equal(t, "Shared Autoscaling", d.Get("cluster_name"))
	assert.Equal(t, "i3.xlarge", d.Get("node_type_id"))
	assert.Equal(t, 4, d.Get("autoscale.0.max_workers"))
	assert.Equal(t, "requests", d.Get("library.2047590238.pypi.0.package"))
	assert.Equal(t, "RUNNING", d.Get("state"))
	assert.Equal(t, false, d.Get("is_pinned"))

	for k, v := range d.State().Attributes {
		fmt.Printf("assert.Equal(t, %#v, d.Get(%#v))\n", v, k)
	}
}

func TestResourceClusterRead_NotFound(t *testing.T) {
	qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/clusters/get?cluster_id=abc",
				Response: common.APIErrorBody{
					ErrorCode: "NOT_FOUND",
					Message:   "Item not found",
				},
				Status: 404,
			},
		},
		Resource: ResourceCluster(),
		Read:     true,
		Removed:  true,
		ID:       "abc",
	}.ApplyNoError(t)
}

func TestResourceClusterRead_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/clusters/get?cluster_id=abc",
				Response: common.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		Resource: ResourceCluster(),
		Read:     true,
		ID:       "abc",
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "abc", d.Id(), "Id should not be empty for error reads")
}

func TestResourceClusterUpdate(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:       "GET",
				Resource:     "/api/2.0/clusters/get?cluster_id=abc",
				ReuseRequest: true,
				Response: ClusterInfo{
					ClusterID:              "abc",
					NumWorkers:             100,
					ClusterName:            "Shared Autoscaling",
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
					AutoterminationMinutes: 15,
					State:                  ClusterStateRunning,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/events",
				ExpectedRequest: EventsRequest{
					ClusterID:  "abc",
					Limit:      1,
					Order:      SortDescending,
					EventTypes: []ClusterEventType{EvTypePinned, EvTypeUnpinned},
				},
				Response: EventsResponse{
					Events:     []ClusterEvent{},
					TotalCount: 0,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/start",
				ExpectedRequest: ClusterID{
					ClusterID: "abc",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{},
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/edit",
				ExpectedRequest: Cluster{
					AutoterminationMinutes: 15,
					ClusterID:              "abc",
					NumWorkers:             100,
					ClusterName:            "Shared Autoscaling",
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{},
				},
			},
		},
		ID:       "abc",
		Update:   true,
		Resource: ResourceCluster(),
		State: map[string]interface{}{
			"autotermination_minutes": 15,
			"cluster_name":            "Shared Autoscaling",
			"spark_version":           "7.1-scala12",
			"node_type_id":            "i3.xlarge",
			"num_workers":             100,
		},
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc", d.Id(), "Id should be the same as in reading")
}

func TestResourceClusterUpdateWithPinned(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:       "GET",
				Resource:     "/api/2.0/clusters/get?cluster_id=abc",
				ReuseRequest: true,
				Response: ClusterInfo{
					ClusterID:              "abc",
					NumWorkers:             100,
					ClusterName:            "Shared Autoscaling",
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
					AutoterminationMinutes: 15,
					State:                  ClusterStateRunning,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/events",
				ExpectedRequest: EventsRequest{
					ClusterID:  "abc",
					Limit:      1,
					Order:      SortDescending,
					EventTypes: []ClusterEventType{EvTypePinned, EvTypeUnpinned},
				},
				Response: EventsResponse{
					Events:     []ClusterEvent{},
					TotalCount: 0,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/start",
				ExpectedRequest: ClusterID{
					ClusterID: "abc",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{},
				},
			},
			{
				Method:          "POST",
				Resource:        "/api/2.0/clusters/pin",
				ExpectedRequest: ClusterID{ClusterID: "abc"},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{},
				},
			},
		},
		ID:       "abc",
		Update:   true,
		Resource: ResourceCluster(),
		InstanceState: map[string]string{
			"autotermination_minutes": "15",
			"cluster_name":            "Shared Autoscaling",
			"spark_version":           "7.1-scala12",
			"node_type_id":            "i3.xlarge",
			"num_workers":             "100",
		},
		State: map[string]interface{}{
			"autotermination_minutes": 15,
			"cluster_name":            "Shared Autoscaling",
			"spark_version":           "7.1-scala12",
			"node_type_id":            "i3.xlarge",
			"num_workers":             100,
			"is_pinned":               true,
		},
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc", d.Id(), "Id should be the same as in reading")
}

func TestResourceClusterUpdate_LibrariesChangeOnTerminatedCluster(t *testing.T) {
	terminated := qa.HTTPFixture{
		Method:   "GET",
		Resource: "/api/2.0/clusters/get?cluster_id=abc",
		Response: ClusterInfo{
			ClusterID:    "abc",
			NumWorkers:   100,
			SparkVersion: "7.1-scala12",
			NodeTypeID:   "i3.xlarge",
			State:        ClusterStateTerminated,
			StateMessage: "Terminated for test reasons",
		},
	}
	newLibs := qa.HTTPFixture{
		Method:   "GET",
		Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
		Response: libraries.ClusterLibraryStatuses{
			ClusterID: "abc",
			LibraryStatuses: []libraries.LibraryStatus{
				{
					Library: &libraries.Library{
						Jar: "dbfs://foo.jar",
					},
					Status: "INSTALLED",
				},
				{
					Library: &libraries.Library{
						Egg: "dbfs://bar.egg",
					},
					Status: "INSTALLED",
				},
			},
		},
	}
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			terminated, // 1 of ...
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/edit",
				ExpectedRequest: Cluster{
					AutoterminationMinutes: 60,
					ClusterID:              "abc",
					NumWorkers:             100,
					SparkVersion:           "7.1-scala12",
					NodeTypeID:             "i3.xlarge",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					ClusterID: "abc",
					LibraryStatuses: []libraries.LibraryStatus{
						{
							Library: &libraries.Library{
								Egg: "dbfs://bar.egg",
							},
							Status: "INSTALLED",
						},
						{
							Library: &libraries.Library{
								Pypi: &libraries.PyPi{
									Package: "requests",
								},
							},
							Status: "INSTALLED",
						},
					},
				},
			},
			{ // check to see if cluster is restarting (if so wait)
				Method:   "GET",
				Resource: "/api/2.0/clusters/get?cluster_id=abc",
				Response: ClusterInfo{
					ClusterID:    "abc",
					NumWorkers:   100,
					SparkVersion: "7.1-scala12",
					NodeTypeID:   "i3.xlarge",
					State:        ClusterStateTerminated,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/events",
				ExpectedRequest: EventsRequest{
					ClusterID:  "abc",
					Limit:      1,
					Order:      SortDescending,
					EventTypes: []ClusterEventType{EvTypePinned, EvTypeUnpinned},
				},
				Response: EventsResponse{
					Events:     []ClusterEvent{},
					TotalCount: 0,
				},
			},
			{ // start cluster before libs install
				Method:   "POST",
				Resource: "/api/2.0/clusters/start",
				ExpectedRequest: ClusterID{
					ClusterID: "abc",
				},
			},
			{ // 2 of ...
				Method:   "GET",
				Resource: "/api/2.0/clusters/get?cluster_id=abc",
				Response: ClusterInfo{
					ClusterID:    "abc",
					NumWorkers:   100,
					SparkVersion: "7.1-scala12",
					NodeTypeID:   "i3.xlarge",
					State:        ClusterStateRunning,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/libraries/uninstall",
				ExpectedRequest: libraries.ClusterLibraryList{
					ClusterID: "abc",
					Libraries: []libraries.Library{
						{
							Pypi: &libraries.PyPi{
								Package: "requests",
							},
						},
					},
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/libraries/install",
				ExpectedRequest: libraries.ClusterLibraryList{
					ClusterID: "abc",
					Libraries: []libraries.Library{
						{
							Jar: "dbfs://foo.jar",
						},
					},
				},
			},
			newLibs,
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/delete",
				ExpectedRequest: ClusterID{
					ClusterID: "abc",
				},
			},
			terminated, // 3 of 4
			// read
			terminated, // 4 of 4
			newLibs,
		},
		ID:       "abc",
		Update:   true,
		Resource: ResourceCluster(),
		HCL: `num_workers = 100
		spark_version = "7.1-scala12"
		node_type_id = "i3.xlarge"

		library {
			jar = "dbfs://foo.jar"
		}

		library {
			egg = "dbfs://bar.egg"
		}`,
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc", d.Id(), "Id should be the same as in reading")
}

func TestResourceClusterUpdate_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "GET",
				Resource: "/api/2.0/clusters/get?cluster_id=abc",
				Response: common.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		ID:       "abc",
		Update:   true,
		Resource: ResourceCluster(),
		State: map[string]interface{}{
			"autotermination_minutes": 15,
			"cluster_name":            "Shared Autoscaling",
			"spark_version":           "7.1-scala12",
			"node_type_id":            "i3.xlarge",
			"num_workers":             100,
		},
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "abc", d.Id())
}

func TestResourceClusterDelete(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/delete",
				ExpectedRequest: map[string]string{
					"cluster_id": "abc",
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/clusters/get?cluster_id=abc",
				Response: ClusterInfo{
					State: ClusterStateTerminated,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/permanent-delete",
				ExpectedRequest: map[string]string{
					"cluster_id": "abc",
				},
			},
		},
		Resource: ResourceCluster(),
		Delete:   true,
		ID:       "abc",
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, "abc", d.Id())
}

func TestResourceClusterDelete_Error(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/delete",
				Response: common.APIErrorBody{
					ErrorCode: "INVALID_REQUEST",
					Message:   "Internal error happened",
				},
				Status: 400,
			},
		},
		Resource: ResourceCluster(),
		Delete:   true,
		ID:       "abc",
	}.Apply(t)
	qa.AssertErrorStartsWith(t, err, "Internal error happened")
	assert.Equal(t, "abc", d.Id())
}

func TestResourceClusterCreate_SingleNode(t *testing.T) {
	d, err := qa.ResourceFixture{
		Fixtures: []qa.HTTPFixture{
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/create",
				ExpectedRequest: Cluster{
					NumWorkers:             0,
					ClusterName:            "Single Node Cluster",
					SparkVersion:           "7.3.x-scala12",
					NodeTypeID:             "Standard_F4s",
					AutoterminationMinutes: 120,
					SparkConf: map[string]string{
						"spark.master":                     "local[*]",
						"spark.databricks.cluster.profile": "singleNode",
					},
					CustomTags: map[string]string{
						"ResourceClass": "SingleNode",
					},
				},
				Response: ClusterInfo{
					ClusterID: "abc",
					State:     ClusterStateRunning,
				},
			},
			{
				Method:   "POST",
				Resource: "/api/2.0/clusters/events",
				ExpectedRequest: EventsRequest{
					ClusterID:  "abc",
					Limit:      1,
					Order:      SortDescending,
					EventTypes: []ClusterEventType{EvTypePinned, EvTypeUnpinned},
				},
				Response: EventsResponse{
					Events:     []ClusterEvent{},
					TotalCount: 0,
				},
			},
			{
				Method:       "GET",
				ReuseRequest: true,
				Resource:     "/api/2.0/clusters/get?cluster_id=abc",
				Response: ClusterInfo{
					ClusterID:              "abc",
					ClusterName:            "Single Node Cluster",
					SparkVersion:           "7.3.x-scala12",
					NodeTypeID:             "Standard_F4s",
					AutoterminationMinutes: 120,
					State:                  ClusterStateRunning,
					SparkConf: map[string]string{
						"spark.master":                     "local[*]",
						"spark.databricks.cluster.profile": "singleNode",
					},
					CustomTags: map[string]string{
						"ResourceClass": "SingleNode",
					},
				},
			},
			{
				Method:   "GET",
				Resource: "/api/2.0/libraries/cluster-status?cluster_id=abc",
				Response: libraries.ClusterLibraryStatuses{
					LibraryStatuses: []libraries.LibraryStatus{},
				},
			},
		},
		Create:   true,
		Resource: ResourceCluster(),
		State: map[string]interface{}{
			"autotermination_minutes": 120,
			"cluster_name":            "Single Node Cluster",
			"spark_version":           "7.3.x-scala12",
			"node_type_id":            "Standard_F4s",
			"is_pinned":               false,
			"spark_conf": map[string]interface{}{
				"spark.master":                     "local[*]",
				"spark.databricks.cluster.profile": "singleNode",
			},
			"custom_tags": map[string]interface{}{
				"ResourceClass": "SingleNode",
			},
		},
	}.Apply(t)
	assert.NoError(t, err, err)
	assert.Equal(t, 0, d.Get("num_workers"))
}

func TestResourceClusterCreate_SingleNodeFail(t *testing.T) {
	_, err := qa.ResourceFixture{
		Create:   true,
		Resource: ResourceCluster(),
		State: map[string]interface{}{
			"autotermination_minutes": 120,
			"cluster_name":            "Single Node Cluster",
			"spark_version":           "7.3.x-scala12",
			"node_type_id":            "Standard_F4s",
			"is_pinned":               false,
		},
	}.Apply(t)
	assert.Error(t, err, err)
	require.Equal(t, true, strings.Contains(err.Error(), "NumWorkers could be 0 only for SingleNode clusters"))
}

func TestResourceClusterCreate_NegativeNumWorkers(t *testing.T) {
	_, err := qa.ResourceFixture{
		Create:   true,
		Resource: ResourceCluster(),
		State: map[string]interface{}{
			"autotermination_minutes": 120,
			"cluster_name":            "Broken Cluster",
			"spark_version":           "7.3.x-scala12",
			"node_type_id":            "Standard_F4s",
			"num_workers":             -10,
		},
	}.Apply(t)
	assert.Error(t, err, err)
	require.Equal(t, true, strings.Contains(err.Error(), "expected num_workers to be at least (0)"))
}

func TestResourceClusterUpdate_FailNumWorkersZero(t *testing.T) {
	_, err := qa.ResourceFixture{
		ID:       "abc",
		Update:   true,
		Resource: ResourceCluster(),
		InstanceState: map[string]string{
			"autotermination_minutes": "15",
			"cluster_name":            "Shared Autoscaling",
			"spark_version":           "7.1-scala12",
			"node_type_id":            "i3.xlarge",
			"num_workers":             "100",
		},
		State: map[string]interface{}{
			"autotermination_minutes": 15,
			"cluster_name":            "Shared Autoscaling",
			"spark_version":           "7.1-scala12",
			"node_type_id":            "i3.xlarge",
			"num_workers":             0,
		},
	}.Apply(t)
	assert.Error(t, err, err)
	require.Equal(t, true, strings.Contains(err.Error(), "NumWorkers could be 0 only for SingleNode clusters"))
}

func TestModifyClusterRequestAws(t *testing.T) {
	c := Cluster{
		InstancePoolID: "a",
		AwsAttributes: &AwsAttributes{
			InstanceProfileArn: "b",
			ZoneID:             "c",
		},
		EnableElasticDisk: true,
		NodeTypeID:        "d",
		DriverNodeTypeID:  "e",
	}
	c.ModifyRequestOnInstancePool()
	assert.Equal(t, "", c.AwsAttributes.ZoneID)
	assert.Equal(t, "", c.NodeTypeID)
	assert.Equal(t, "", c.DriverNodeTypeID)
	assert.Equal(t, false, c.EnableElasticDisk)
}

func TestModifyClusterRequestAzure(t *testing.T) {
	c := Cluster{
		InstancePoolID: "a",
		AzureAttributes: &AzureAttributes{
			FirstOnDemand: 1,
		},
		EnableElasticDisk: true,
		NodeTypeID:        "d",
		DriverNodeTypeID:  "e",
	}
	c.ModifyRequestOnInstancePool()
	assert.Nil(t, c.AzureAttributes)
	assert.Equal(t, "", c.NodeTypeID)
	assert.Equal(t, "", c.DriverNodeTypeID)
	assert.Equal(t, false, c.EnableElasticDisk)
}

func TestModifyClusterRequestGcp(t *testing.T) {
	c := Cluster{
		InstancePoolID: "a",
		GcpAttributes: &GcpAttributes{
			UsePreemptibleExecutors: true,
		},
		EnableElasticDisk: true,
		NodeTypeID:        "d",
		DriverNodeTypeID:  "e",
	}
	c.ModifyRequestOnInstancePool()
	assert.Equal(t, false, c.GcpAttributes.UsePreemptibleExecutors)
	assert.Equal(t, "", c.NodeTypeID)
	assert.Equal(t, "", c.DriverNodeTypeID)
	assert.Equal(t, false, c.EnableElasticDisk)
}

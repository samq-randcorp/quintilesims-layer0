package logic

import (
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

func TestGetEnvironment(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retEnvironment := &models.Environment{
		EnvironmentID:   "e1",
		ClusterCount:    1,
		InstanceSize:    "m3.medium",
		SecurityGroupID: "sgid",
	}

	testLogic.Backend.EXPECT().
		GetEnvironment("e1").
		Return(retEnvironment, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env"},
		{EntityID: "extra", EntityType: "environment", Key: "name", Value: "extra"},
	})

	environmentLogic := NewL0EnvironmentLogic(testLogic.Logic())
	received, err := environmentLogic.GetEnvironment("e1")
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Environment{
		EnvironmentID:   "e1",
		EnvironmentName: "env",
		ClusterCount:    1,
		InstanceSize:    "m3.medium",
		SecurityGroupID: "sgid",
	}

	testutils.AssertEqual(t, received, expected)
}

func TestListEnvironments(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retEnvironments := []*models.Environment{
		{
			EnvironmentID:   "e1",
			ClusterCount:    1,
			InstanceSize:    "m3.medium",
			SecurityGroupID: "sgid1",
		},
		{
			EnvironmentID:   "e2",
			ClusterCount:    2,
			InstanceSize:    "m3.small",
			SecurityGroupID: "sgid2",
		},
	}

	testLogic.Backend.EXPECT().
		ListEnvironments().
		Return(retEnvironments, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env_1"},
		{EntityID: "e2", EntityType: "environment", Key: "name", Value: "env_2"},
		{EntityID: "extra", EntityType: "environment", Key: "name", Value: "extra"},
	})

	environmentLogic := NewL0EnvironmentLogic(testLogic.Logic())
	environments, err := environmentLogic.ListEnvironments()
	if err != nil {
		t.Fatal(err)
	}

	expected := []*models.Environment{
		{
			EnvironmentID:   "e1",
			EnvironmentName: "env_1",
			ClusterCount:    1,
			InstanceSize:    "m3.medium",
			SecurityGroupID: "sgid1",
		},
		{
			EnvironmentID:   "e2",
			EnvironmentName: "env_2",
			ClusterCount:    2,
			InstanceSize:    "m3.small",
			SecurityGroupID: "sgid2",
		},
	}

	testutils.AssertEqual(t, len(environments), 2)
	testutils.AssertEqual(t, environments[0], expected[0])
	testutils.AssertEqual(t, environments[1], expected[1])
}

func TestDeleteEnvironment(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.Backend.EXPECT().
		DeleteEnvironment("e1").
		Return(nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env"},
		{EntityID: "extra", EntityType: "environment", Key: "name", Value: "extra"},
	})

	environmentLogic := NewL0EnvironmentLogic(testLogic.Logic())
	if err := environmentLogic.DeleteEnvironment("e1"); err != nil {
		t.Fatal(err)
	}

	tags, err := testLogic.TagStore.SelectAll()
	if err != nil {
		t.Fatal(err)
	}

	// make sure the 'extra' tag is the only one left
	testutils.AssertEqual(t, len(tags), 1)
}

func TestCanCreateEnvironment(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env_1"},
		{EntityID: "e2", EntityType: "environment", Key: "name", Value: "env_2"},
		{EntityID: "extra", EntityType: "environment", Key: "name", Value: "extra"},
	})

	environmentLogic := NewL0EnvironmentLogic(testLogic.Logic())

	cases := map[string]bool{
		"env_1":  false,
		"env_2":  false,
		"env3":   true,
		"env_12": true,
		"env":    true,
	}

	for name, expected := range cases {
		request := models.CreateEnvironmentRequest{EnvironmentName: name}

		received, err := environmentLogic.CanCreateEnvironment(request)
		if err != nil {
			t.Fatal(err)
		}

		if received != expected {
			t.Errorf("Failure on case '%s': response was %v, expected %v", name, received, expected)
		}
	}
}

func TestCreateEnvironment(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retEnvironment := &models.Environment{
		EnvironmentID:   "e1",
		ClusterCount:    1,
		InstanceSize:    "m3.medium",
		SecurityGroupID: "sgid",
	}

	testLogic.Backend.EXPECT().
		CreateEnvironment("name", "m3.medium", 2, []byte("user_data")).
		Return(retEnvironment, nil)

	request := models.CreateEnvironmentRequest{
		EnvironmentName:  "name",
		InstanceSize:     "m3.medium",
		MinClusterCount:  2,
		UserDataTemplate: []byte("user_data"),
	}

	environmentLogic := NewL0EnvironmentLogic(testLogic.Logic())
	received, err := environmentLogic.CreateEnvironment(request)
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Environment{
		EnvironmentID:   "e1",
		EnvironmentName: "name",
		ClusterCount:    1,
		InstanceSize:    "m3.medium",
		SecurityGroupID: "sgid",
	}

	testutils.AssertEqual(t, received, expected)
	testLogic.AssertTagExists(t, models.Tag{EntityID: "e1", EntityType: "environment", Key: "name", Value: "name"})
}

func TestUpdateEnvironment(t *testing.T) {
	testLogic, ctrl := NewTestLogic(t)
	defer ctrl.Finish()

	retEnvironment := &models.Environment{
		EnvironmentID:   "e1",
		ClusterCount:    2,
		InstanceSize:    "m3.medium",
		SecurityGroupID: "sgid",
	}

	testLogic.Backend.EXPECT().
		UpdateEnvironment("e1", 2).
		Return(retEnvironment, nil)

	testLogic.AddTags(t, []*models.Tag{
		{EntityID: "e1", EntityType: "environment", Key: "name", Value: "env"},
		{EntityID: "extra", EntityType: "environment", Key: "name", Value: "extra"},
	})

	environmentLogic := NewL0EnvironmentLogic(testLogic.Logic())
	received, err := environmentLogic.UpdateEnvironment("e1", 2)
	if err != nil {
		t.Fatal(err)
	}

	expected := &models.Environment{
		EnvironmentID:   "e1",
		EnvironmentName: "env",
		ClusterCount:    2,
		InstanceSize:    "m3.medium",
		SecurityGroupID: "sgid",
	}

	testutils.AssertEqual(t, received, expected)
}

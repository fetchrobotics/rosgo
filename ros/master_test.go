package ros

import (
	"os"
	"testing"
)

func TestCallRosAPI(t *testing.T) {
	masterURI := os.Getenv("ROS_MASTER_URI")
	nodeName := "/test_call_ros_api"

	paramName := "/test_param"
	nonExistingParam := "/test_param_unknown"
	paramValue := int32(100)

	if _, err := callRosAPI(masterURI, "setParam", nodeName, paramName, paramValue); err != nil {
		t.Errorf("Error setting ROS param using callRosAPI: %v", err)
	}

	if result, err := callRosAPI(masterURI, "searchParam", nodeName, paramName); err != nil {
		t.Errorf("Error searching ROS param using callRosAPI: %v", err)
	} else {
		found, ok := result.(string)
		if !ok {
			t.Errorf("Error asserting searchParam result to string")
		}
		if found == "false" {
			t.Errorf("Param %s not found after setting the param", paramName)
		}
	}

	if result, err := callRosAPI(masterURI, "hasParam", nodeName, paramName); err != nil {
		t.Errorf("Error checking if param exists using callRosAPI: %v", err)
	} else {
		hasParam, ok := result.(bool)
		if !ok {
			t.Errorf("Error asserting hasParam result to boolean")
		}
		if !hasParam {
			t.Errorf("Param %s not found after setting the param", paramName)
		}
	}

	if result, err := callRosAPI(masterURI, "getParam", nodeName, paramName); err != nil {
		t.Errorf("Error getting param using callRosAPI: %v", err)
	} else {
		paramGot, ok := result.(int32)
		if !ok {
			t.Errorf("Error asserting getParam result to int")
		}
		if paramValue != paramGot {
			t.Errorf("Param value received is not equal to param value set: Wanted %v Got %v",
				paramValue, paramGot)
		}
	}

	if _, err := callRosAPI(masterURI, "getParam", nodeName, nonExistingParam); err == nil {
		t.Errorf("Expected API call to fail")
	}
}

func TestBuildRosAPIResult(t *testing.T) {
	var statusWanted int32 = errorStatus
	messageWanted := "Not implemented"
	valueWanted := 0

	result := buildRosAPIResult(statusWanted, messageWanted, valueWanted)
	if result == nil {
		t.Errorf("Nil result received")
	}

	res, ok := result.([]interface{})
	if !ok {
		t.Errorf("Result is malformed: Result should be an array of interface values.")
	}

	if len(res) != 3 {
		t.Errorf("Result is malformed: Length of result array should be 3")
	}

	statusGot, ok := res[0].(int32)
	if !ok {
		t.Errorf("Error getting status code from result. Type assertion failed.")
		return
	}

	messageGot, ok := res[1].(string)
	if !ok {
		t.Errorf("Error getting message from result. Type assertion failed.")
		return
	}
	valueGot, ok := res[2].(int)
	if !ok {
		t.Errorf("Error getting message from result. Type assertion failed.")
		return
	}

	if statusWanted != statusGot {
		t.Errorf("Expected %v but got %v", statusWanted, statusGot)
	}
	if messageWanted != messageGot {
		t.Errorf("Expected %v but got %v", messageWanted, messageGot)
	}
	if valueWanted != valueGot {
		t.Errorf("Expected %v but got %v", valueWanted, valueGot)
	}
}

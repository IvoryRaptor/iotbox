package helper
import "fmt"

func cleanupInterfaceArray(in []interface{}) []interface{} {
    res := make([]interface{}, len(in))
    for i, v := range in {
        res[i] = CleanupMapValue(v)
    }
    return res
}

func cleanupInterfaceMap(in map[interface{}]interface{}) map[string]interface{} {
    res := make(map[string]interface{})
    for k, v := range in {
        res[fmt.Sprintf("%v", k)] = CleanupMapValue(v)
    }
    return res
}
// CleanupMapValue 清理map[interface{}]interface{} 转成map[string]interface{}
func CleanupMapValue(v interface{}) interface{} {
    switch v := v.(type) {
    case []interface{}:
        return cleanupInterfaceArray(v)
    case map[interface{}]interface{}:
		return cleanupInterfaceMap(v)
    default:
        return v
    }
}
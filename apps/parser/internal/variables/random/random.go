package random

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
	types "tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

const exampleStr = "Please check example."
var Variable = types.Variable{
	Name:        "random",
	Description: lo.ToPtr("Random number from N to N"),
	Example:     lo.ToPtr("random|1-100"),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}
		rand.Seed(time.Now().UnixNano())

		params := [2]int{}
		if data.Params == nil {
			result.Result = "You have not passed params to random variable. " + exampleStr
			return result, nil
		}

		if len(*data.Params) != 2 {
			result.Result = "Wrong number of arguments passed to random. " + exampleStr
			return result, nil
		}

		if data.Params != nil {
			parsed := strings.Split(*data.Params, "-")

			if len(parsed) < 2 {
				result.Result = "you have no passed 2 params for random. " + exampleStr
				return result, nil
			}

			first, err := strconv.Atoi(parsed[0])
			if err == nil {
				params[0] = first
			} else {
				result.Result = "cannot parse first number from arguments. " + exampleStr
				return result, nil
			}
			second, err := strconv.Atoi(parsed[1])
			if err == nil {
				params[1] = second
			} else {
				result.Result = "cannot parse second number from arguments. " + exampleStr
				return result, nil
			}
		}

		if params[0] > params[1] {
			return nil, errors.New("first number cannot be larger then second")
		}

		if params[0] < 0 || params[1] < 0 {
			return nil, errors.New("numbers cannot be lower then 0")
		}

		random := params[0] + rand.Intn(params[1]-params[0]+1)
		result.Result = strconv.Itoa(random)

		return result, nil
	},
}

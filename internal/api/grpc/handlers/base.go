package handlers

import (
	"fmt"
	"strconv"

	"github.com/ose-micro/core/dto"
	ose_error "github.com/ose-micro/error"
	commonv1 "github.com/ose-micro/postman/internal/api/grpc/gen/go/ose/micro/common/v1"
)

func convertStringMapToInterfaceMap(in map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func stringifyInterfaceMap(in map[string]interface{}) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out
}

func buildAppRequest(query *commonv1.Request) (*dto.Request, error) {
	if query == nil {
		return nil, fmt.Errorf("query is nil")
	}

	facets := make([]dto.Query, len(query.Facets))

	for i, facet := range query.Facets {

		filters := make([]dto.Filter, len(facet.Filters))

		for i, filter := range facet.Filters {
			filters[i] = dto.Filter{
				Field: filter.Field,
				Op: func() dto.FilterOp {
					switch filter.Op {
					case commonv1.FilterOp_FILTER_OP_EQ:
						return dto.OpEq
					case commonv1.FilterOp_FILTER_OP_GTE:
						return dto.OpGte
					case commonv1.FilterOp_FILTER_OP_GT:
						return dto.OpGt
					case commonv1.FilterOp_FILTER_OP_LT:
						return dto.OpLt
					case commonv1.FilterOp_FILTER_OP_LTE:
						return dto.OpLte
					case commonv1.FilterOp_FILTER_OP_IN:
						return dto.OpIn
					case commonv1.FilterOp_FILTER_OP_NE:
						return dto.OpNe
					case commonv1.FilterOp_FILTER_OP_NIN:
						return dto.OpNin
					default:
						return dto.OpEq
					}
				}(),
				Value: func() interface{} {
					switch filter.Op {
					case commonv1.FilterOp_FILTER_OP_EQ, commonv1.FilterOp_FILTER_OP_IN, commonv1.FilterOp_FILTER_OP_NE,
						commonv1.FilterOp_FILTER_OP_NIN:
						return filter.Value
					case commonv1.FilterOp_FILTER_OP_GTE, commonv1.FilterOp_FILTER_OP_GT, commonv1.FilterOp_FILTER_OP_LT,
						commonv1.FilterOp_FILTER_OP_LTE:
						value, err := strconv.ParseFloat(filter.Value, 32)
						if err != nil {
							return nil
						}
						return float32(value)

					default:
						return dto.OpEq
					}
				}(),
			}
		}

		aggregations := make([]dto.Aggregation, len(facet.Aggregations))
		for i, agg := range facet.Aggregations {
			aggregations[i] = dto.Aggregation{
				Field: agg.Field,
				As:    agg.As,
				Type: func() dto.AggregationType {
					switch agg.Type {
					case commonv1.AggregationType_AGGREGATION_TYPE_SUM:
						return dto.AggSum
					case commonv1.AggregationType_AGGREGATION_TYPE_AVG:
						return dto.AggAvg
					case commonv1.AggregationType_AGGREGATION_TYPE_COUNT:
						return dto.AggCount
					case commonv1.AggregationType_AGGREGATION_TYPE_MAX:
						return dto.AggMax
					case commonv1.AggregationType_AGGREGATION_TYPE_MIN:
						return dto.AggMin
					default:
						return dto.AggCount
					}
				}(),
			}
		}

		sorts := make([]dto.SortOption, len(facet.Sort))
		for i, sort := range facet.Sort {
			sorts[i] = dto.SortOption{
				Field: sort.Field,
				Order: func() dto.SortOrder {
					switch sort.Order {
					case commonv1.SortOrder_SORT_ORDER_ASC:
						return dto.SortAsc
					case commonv1.SortOrder_SORT_ORDER_DESC:
						return dto.SortDesc
					default:
						return dto.SortAsc
					}
				}(),
			}
		}

		computedFields := make([]dto.ComputedField, len(facet.ComputedFields))
		for i, computed := range facet.ComputedFields {
			computedFields[i] = dto.ComputedField{
				Name: computed.Name,
				Operator: func() dto.ComputedOperator {
					switch computed.Operator {
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_ADD:
						return dto.OpAdd
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_SUBTRACT:
						return dto.OpSubtract
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_MULTIPLY:
						return dto.OpMultiply
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_DIVIDE:
						return dto.OpDivide
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_CONCAT:
						return dto.OpConcat
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_DATE_TRUNC:
						return dto.OpDateTrunc
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_DAY_OF_MONTH:
						return dto.OpDateTrunc
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_IF_NULL:
						return dto.OpIfNull
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_MONTH:
						return dto.OpDayOfMonth
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_SUBSTR:
						return dto.OpSubstr
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_TO_LOWER:
						return dto.OpToLower
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_TO_UPPER:
						return dto.OpToUpper
					case commonv1.ComputedOperator_COMPUTED_OPERATOR_YEAR:
						return dto.OpYear
					default:
						return dto.OpAdd
					}
				}(),
				Operands: computed.Operands,
			}
		}

		facets[i] = dto.Query{
			Name:           facet.Name,
			Filters:        filters,
			GroupBy:        facet.GroupBy,
			Aggregations:   aggregations,
			Sort:           sorts,
			ComputedFields: computedFields,
			Skip:           facet.Skip,
			Limit:          facet.Limit,
		}
	}

	return &dto.Request{
		Queries: facets,
	}, nil
}

func parseError(err error) error {
	oseErr := ose_error.CastToError(err)

	err = ose_error.ToGRPCError(oseErr)
	if err != nil {
		return err
	}

	return nil
}

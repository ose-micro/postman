package handlers

import (
	"fmt"
	"time"

	commonv1 "github.com/moriba-cloud/ose-postman/internal/interface/grpc/gen/go/common/v1"
	"github.com/ose-micro/core/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	TABLE string
)

func operatorToEnum(filter dto.Operator) commonv1.Operator {
	switch filter {
	case dto.EQUAL:
		return commonv1.Operator_EQUAL
	case dto.LIKE:
		return commonv1.Operator_LIKE
	case dto.BETWEEN:
		return commonv1.Operator_BETWEEN
	case dto.AFTER:
		return commonv1.Operator_AFTER
	case dto.BEFORE:
		return commonv1.Operator_BEFORE
	case dto.GREATER_THAN:
		return commonv1.Operator_GREATER_THAN
	case dto.LESS_THAN:
		return commonv1.Operator_LESS_THAN
	case dto.GREATER_THAN_EQUAL:
		return commonv1.Operator_GREATER_THAN_EQUAL
	case dto.LESS_THAN_EQUAL:
		return commonv1.Operator_LESS_THAN_EQUAL
	case dto.DATE_EQUAL:
		return commonv1.Operator_DATE_EQUAL
	case dto.DATE_BETWEEN:
		return commonv1.Operator_DATE_BETWEEN
	default:
		return commonv1.Operator_UNKNOWN_FILTER
	}
}

func orderToEnum(sort dto.Direction) commonv1.Order {
	switch sort {
	case dto.ASC:
		return commonv1.Order_ASC
	case dto.DESC:
		return commonv1.Order_DESC
	default:
		return commonv1.Order_ASC
	}
}

func enumToOperator(filter commonv1.Operator) dto.Operator {
	switch filter {
	case commonv1.Operator_EQUAL:
		return dto.EQUAL
	case commonv1.Operator_LIKE:
		return dto.LIKE
	case commonv1.Operator_BETWEEN:
		return dto.BETWEEN
	case commonv1.Operator_AFTER:
		return dto.AFTER
	case commonv1.Operator_BEFORE:
		return dto.BEFORE
	case commonv1.Operator_GREATER_THAN:
		return dto.GREATER_THAN
	case commonv1.Operator_LESS_THAN:
		return dto.LESS_THAN
	case commonv1.Operator_GREATER_THAN_EQUAL:
		return dto.GREATER_THAN_EQUAL
	case commonv1.Operator_LESS_THAN_EQUAL:
		return dto.LESS_THAN_EQUAL
	case commonv1.Operator_DATE_EQUAL:
		return dto.DATE_EQUAL
	case commonv1.Operator_DATE_BETWEEN:
		return dto.DATE_BETWEEN
	default:
		return ""
	}
}

func enumToOrder(sort commonv1.Order) dto.Direction {
	switch sort {
	case commonv1.Order_ASC:
		return dto.ASC
	case commonv1.Order_DESC:
		return dto.DESC
	default:
		return dto.ASC
	}
}

func processValue(filter *commonv1.Filter) interface{} {
	switch v := filter.Value.(type) {
	case *commonv1.Filter_StringValue:
		return v.StringValue
	case *commonv1.Filter_Int32Value:
		return v.Int32Value
	case *commonv1.Filter_TimeValue:
		return v.TimeValue.AsTime()
	case *commonv1.Filter_ValuesValue:
		return v.ValuesValue.Values
	default:
		return ""
	}
}

func filterToGrpc(filter dto.Filter) *commonv1.Filter {
	operator := operatorToEnum(filter.Operator)
	field := filter.Field

	switch v := filter.Value.(type) {
	case string:
		return &commonv1.Filter{
			Field:      field,
			Operator: operator,
			Value:    &commonv1.Filter_StringValue{StringValue: v},
		}
	case int32:
		return &commonv1.Filter{
			Field:      field,
			Operator: operator,
			Value:    &commonv1.Filter_Int32Value{Int32Value: v},
		}
	case time.Time:
		return &commonv1.Filter{
			Field:      field,
			Operator: operator,
			Value:    &commonv1.Filter_TimeValue{TimeValue: timestamppb.New(v)},
		}
	case []string:
		return &commonv1.Filter{
			Field:      field,
			Operator: operator,
			Value:    &commonv1.Filter_ValuesValue{ValuesValue: &commonv1.Values{Values: v}},
		}
	default:
		return nil
	}
}

func buildGRPCRequest(request *dto.Request) (*commonv1.Request, error) {
	filters := make([]*commonv1.Filter, len(request.Filter))
	sorts := make([]*commonv1.Sort, len(request.Sort))

	for i, filter := range request.Filter {
		filters[i] = filterToGrpc(filter)
	}

	for i, sort := range request.Sort {
		sorts[i] = &commonv1.Sort{
			Field: sort.Field,
			Order: orderToEnum(sort.Value),
		}
	}

	return &commonv1.Request{
		Pagination: &commonv1.Pagination{
			Page:  int32(request.Pagination.Page),
			Limit: int32(request.Pagination.Limit),
		},
		Filters: filters,
		Sort:    sorts,
	}, nil
}

func buildAppRequest(query *commonv1.Request) (*dto.Request, error) {
	if query == nil {
		return nil, fmt.Errorf("query is nil")
	}

	filters := make([]dto.Filter, len(query.Filters))
	sorts := make([]dto.Sort, len(query.Sort))

	for i, filter := range query.Filters {
		filters[i] = dto.Filter{
			Field:      filter.Field,
			Operator: enumToOperator(filter.Operator),
			Value:    processValue(filter),
		}
	}

	for i, sort := range query.Sort {
		sorts[i] = dto.Sort{
			Field:   sort.Field,
			Value: enumToOrder(sort.Order),
		}
	}

	return &dto.Request{
		Pagination: &dto.Pagination{
			Page: int(query.Pagination.Page),
			Limit: int(query.Pagination.Limit),
		},
		Filter: filters,
		Sort: sorts,
	}, nil
}

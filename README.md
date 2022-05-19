# pagination

Provides generic support for pagination towards a datasource, only SQL is supported out of the box but it is easy to provide additional impl since it's based on this simple interface:

```go
type Page interface {
    Token() string
    Size() int
    Offset() int
    NextToken(resultSize int) *string
}
```

This is based on the [list pagination design pattern](https://cloud.google.com/apis/design/design_patterns#list_pagination) suggested by Google when working with protocol buffers as well as [gorm scopes](https://gorm.io/docs/scopes.html#Pagination).

## Install

```
go get github.com/dentech-floss/pagination@v0.1.0
```

## Usage

So based on the mentioned design pattern above, we have designed our gRPC api like this to provide support for pagination:

```proto
message FindAppointmentsRequest {
  repeated string clinic_ids = 1 [json_name = "clinic_id"];

  google.protobuf.StringValue page_token = 11 [json_name = "page_token"];
  google.protobuf.Int32Value page_size = 12 [json_name = "page_size"];
}
```

```proto
message FindAppointmentsResponse {
  repeated AppointmentDTO appointments = 1;

  string page_token = 11 [json_name = "page_token"];
  int32 page_size = 12 [json_name = "page_size"];
  google.protobuf.StringValue next_page_token = 13 [json_name = "next_page_token"];
}
```

Then in the gRPC server we create a SQL page that we pass to the repository tier:

```go
package example

import (
    "github.com/dentech-floss/pagination/pkg/pagination"

    patient_gateway_service_v1 "go.buf.build/dentechse/go-grpc-gateway-openapiv2/dentechse/patient-api-gateway/api/patient/v1"
)

const (
	DEFAULT_FIND_APPOINTMENTS_PAGE_SIZE = 100
	MAX_FIND_APPOINTMENTS_PAGE_SIZE     = 1000
)

func (s *PatientGatewayServiceV1) FindAppointments(
	ctx context.Context,
	request *patient_gateway_service_v1.FindAppointmentsRequest,
) (*patient_gateway_service_v1.FindAppointmentsResponse, error) {

    var pageToken *string = nil
    if request.PageToken != nil {
        pageToken = &request.PageToken.Value
    }

    var pageSize *int = nil
    if request.PageSize != nil {
        tmp := int(request.PageSize.Value)
        pageSize = &tmp
    }

    page, err := pagination.NewSqlPage(
        pageToken,
        pageSize,
        DEFAULT_FIND_APPOINTMENTS_PAGE_SIZE,
        MAX_FIND_APPOINTMENTS_PAGE_SIZE,
    )
    if err != nil {
        // handle the error
    }

    appointments, err := s.repo.FindAppointmentsForClinics(ctx, clinicIds, page)
    if err != nil {
        // handle the error
    }

    return &patient_gateway_service_v1.FindAppointmentsResponse{
        Appointments: s.appointmentsToDTOs(appointments),
        PageToken:     page.Token(),
        PageSize:      int32(page.Size()),
        NextPageToken: util.StringToWrapper(page.NextToken(len(appointments))),
    }
}
```

The repository tier along with a GORM/SQL implementation looks something like this:

```go
package example

import (
    "github.com/dentech-floss/pagination/pkg/pagination"
)

type Repository interface {
    FindAppointmentsForClinics(ctx context.Context, clinicIds []int32, page pagination.Page) ([]*model.Appointment, error)
}
```

```go
package example

import (
    "github.com/dentech-floss/pagination/pkg/pagination"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

func (r *sqlRepository) FindAppointmentsForClinics(
    ctx context.Context,
    clinicIds []int32,
    page pagination.Page,
) ([]*model.Appointment, error) {

    appointments := make([]*model.Appointment, 0)
    if err := r.db.
        WithContext(ctx). // to propagate the active span for tracing
        Where("clinic_id IN ?", clinicIds).
        Order("start_time asc").
        Scopes(paginationScope(page)).
        Find(&appointments).Error; err != nil {
        return nil, err
    }
    return appointments, nil
}

func paginationScope(page pagination.Page) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Offset(page.Offset()).Limit(page.Size())
    }
}
```
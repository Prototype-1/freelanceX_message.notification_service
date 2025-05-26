package client

import (
	"context"
    pb "github.com/Prototype-1/freelanceX_message.notification_service/proto/invoice_service"
)

type InvoiceServiceClient struct {
    Grpc pb.InvoiceServiceClient
}

func NewInvoiceServiceClient(grpcClient pb.InvoiceServiceClient) *InvoiceServiceClient {
    return &InvoiceServiceClient{Grpc: grpcClient}
}

func (ic *InvoiceServiceClient) GetInvoicePDF(invoiceID string) ([]byte, error) {
    resp, err := ic.Grpc.GetInvoicePDF(context.Background(), &pb.GetInvoicePDFRequest{
        InvoiceId: invoiceID,
    })
    if err != nil {
        return nil, err
    }
    return resp.PdfData, nil
}

package tests

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang_restapi/internal/transaction/config"
	"golang_restapi/internal/transaction/model"
	"golang_restapi/internal/transaction/service"
	"golang_restapi/internal/transaction/service/mocks"
	"testing"
	"time"
)

type TestEntity struct {
	name           string
	userID         uint64
	amount         float64
	setupMocks     func(repository *mocks.MockRepository, publisher *mocks.MockKafkaPublisher, config *config.Config)
	expectedError  bool
	expectedResult model.TransactionDetails
}

func TestTransactionService_Deposit(t *testing.T) {
	tests := []TestEntity{
		{
			name:   "successful deposit",
			userID: 1,
			amount: 100.50,
			setupMocks: func(mockRepo *mocks.MockRepository, mockKafka *mocks.MockKafkaPublisher, config *config.Config) {
				expectedParams := model.DepositParams{
					UserID: 1,
					Amount: 100.50,
				}
				expectedResult := model.TransactionDetails{
					Transaction: model.Transaction{
						ID:        1,
						UserID:    1,
						Amount:    10050,
						Status:    model.TransactionStatusPending,
						Type:      model.TransactionTypeDeposit,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Entries: []model.TransactionEntry{},
				}
				mockRepo.EXPECT().
					Deposit(gomock.Any(), expectedParams).
					Return(expectedResult, nil)
			},
			expectedError: false,
			expectedResult: model.TransactionDetails{
				Transaction: model.Transaction{
					ID:        1,
					UserID:    1,
					Amount:    10050,
					Status:    model.TransactionStatusPending,
					Type:      model.TransactionTypeDeposit,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Entries: []model.TransactionEntry{},
			},
		},
		{
			name:   "repository error",
			userID: 1,
			amount: 100.50,
			setupMocks: func(mockRepo *mocks.MockRepository, publisher *mocks.MockKafkaPublisher, config *config.Config) {
				expectedParams := model.DepositParams{
					UserID: 1,
					Amount: 100.50,
				}
				mockRepo.EXPECT().
					Deposit(gomock.Any(), expectedParams).
					Return(model.TransactionDetails{}, errors.New("database error"))
			},
			expectedError:  true,
			expectedResult: model.TransactionDetails{},
		},

		{
			name:   "kafka publisher error",
			userID: 1,
			amount: 100.50,
			setupMocks: func(mockRepo *mocks.MockRepository, mockKafka *mocks.MockKafkaPublisher, config *config.Config) {
				expectedParams := model.DepositParams{
					UserID: 1,
					Amount: 100.50,
				}
				expectedResult := model.TransactionDetails{
					Transaction: model.Transaction{
						ID:        1,
						UserID:    1,
						Amount:    10050,
						Status:    model.TransactionStatusPending,
						Type:      model.TransactionTypeDeposit,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					Entries: []model.TransactionEntry{},
				}

				mockRepo.EXPECT().
					Deposit(gomock.Any(), expectedParams).
					Return(expectedResult, nil)
				mockKafka.EXPECT().
					Publish(gomock.Any(), config.KafkaPublishTopic, "1", gomock.Any()).
					Return(errors.New("kafka publish error"))

				mockRepo.EXPECT().
					UpdateTransactionStatus(gomock.Any(), uint64(1), model.TransactionStatusFailed).
					Return(nil)
			},
			expectedError:  true,
			expectedResult: model.TransactionDetails{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cfg, err := config.Load()
			if err != nil {
				t.Fatal(err)
			}

			mockRepo := mocks.NewMockRepository(ctrl)
			mockAccountService := mocks.NewMockAccountService(ctrl)
			mockKafka := mocks.NewMockKafkaPublisher(ctrl)
			logger := zerolog.Nop()

			tt.setupMocks(mockRepo, mockKafka, cfg)
			service := service.New(mockRepo, mockAccountService, mockKafka, cfg, &logger)

			result, err := service.Deposit(context.Background(), tt.userID, tt.amount)
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Transaction.ID, result.Transaction.ID)
				assert.Equal(t, tt.expectedResult.Transaction.UserID, result.Transaction.UserID)
				assert.Equal(t, tt.expectedResult.Transaction.Amount, result.Transaction.Amount)
				assert.Equal(t, tt.expectedResult.Transaction.Status, result.Transaction.Status)
				assert.Equal(t, tt.expectedResult.Transaction.Type, result.Transaction.Type)
			}
		})
	}
}

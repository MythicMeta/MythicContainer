package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MythicMeta/MythicContainer/grpc/services"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/translationstructs"
	"github.com/MythicMeta/MythicContainer/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"sync"
	"time"
)

const grpcReconnectDelay = time.Duration(5)

func Initialize(translationContainerName string) {
	opts := []grpc.DialOption{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	connectionString := fmt.Sprintf("%s:%d", utils.MythicConfig.MythicServerHost, utils.MythicConfig.MythicServerGRPCPort)
	for {
		logging.LogDebug("Attempting to connect to grpc...")
		if conn, err := grpc.Dial(connectionString, opts...); err != nil {
			logging.LogError(err, "Failed to connect to GRPC port for Mythic", "connection", connectionString)

		} else {
			var wg sync.WaitGroup
			logging.LogInfo("Successfully connected to grpc", "connection", connectionString)
			client := services.NewTranslationContainerClient(conn)
			closedConnection := make(chan bool, 3)
			wg.Add(3)
			go handleGenerateKeys(&wg, translationContainerName, client, closedConnection)
			go handleMythicToCustomFormat(&wg, translationContainerName, client, closedConnection)
			go handleCustomToMythicFormat(&wg, translationContainerName, client, closedConnection)
			<-closedConnection
			logging.LogInfo("Lost connection to grpc, waiting for all grpc functions to exit...")
			conn.Close()
			wg.Wait()
			close(closedConnection) // use this as a signal that the others should close as well

			logging.LogDebug("All grpc connections closed, opening new ones...")
		}
		time.Sleep(grpcReconnectDelay * time.Second)
	}

}

func handleGenerateKeys(wg *sync.WaitGroup, translationContainerName string, client services.TranslationContainerClient, closedConnection chan bool) {
	defer wg.Done()
	if stream, err := client.GenerateEncryptionKeys(context.Background()); err != nil {
		logging.LogError(err, "Failed to connect to grpc for generate encryption keys")
		closedConnection <- true
		return
	} else if err := stream.Send(&services.TrGenerateEncryptionKeysMessageResponse{
		TranslationContainerName: translationContainerName,
	}); err != nil {
		logging.LogError(err, "Failed to send message to grpc stream")
		closedConnection <- true
		return
	} else {
		for {
			//logging.LogDebug("in for loop waiting to receive from grpc generate encryption keys")
			if input, err := stream.Recv(); err == io.EOF {
				logging.LogError(err, "got EOF from other side")
				closedConnection <- true
				return
			} else if err != nil {
				logging.LogError(err, "Failed to read from stream for generating encryption keys")
				closedConnection <- true
				return
			} else {
				if translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().GenerateEncryptionKeys != nil {
					response := translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().GenerateEncryptionKeys(translationstructs.TrGenerateEncryptionKeysMessage{
						TranslationContainerName: input.GetTranslationContainerName(),
						C2Name:                   input.GetC2Name(),
						CryptoParamValue:         input.GetCryptoParamValue(),
						CryptoParamName:          input.GetCryptoParamName(),
					})
					sendResp := services.TrGenerateEncryptionKeysMessageResponse{
						Success:                  response.Success,
						Error:                    response.Error,
						DecryptionKey:            nil,
						TranslationContainerName: input.GetTranslationContainerName(),
					}
					if response.EncryptionKey != nil {
						sendResp.EncryptionKey = *response.EncryptionKey
					}
					if response.DecryptionKey != nil {
						sendResp.DecryptionKey = *response.DecryptionKey
					}
					if err := stream.Send(&sendResp); err != nil {
						logging.LogError(err, "Failed to send response back to Mythic over grpc")
					}
				}
			}
		}
	}
}
func handleCustomToMythicFormat(wg *sync.WaitGroup, translationContainerName string, client services.TranslationContainerClient, closedConnection chan bool) {
	defer wg.Done()
	if stream, err := client.TranslateFromCustomToMythicFormat(context.Background()); err != nil {
		logging.LogError(err, "Failed to connect to grpc for TranslateFromCustomToMythicFormat")
		closedConnection <- true
		return
	} else if err := stream.Send(&services.TrCustomMessageToMythicC2FormatMessageResponse{
		TranslationContainerName: translationContainerName,
	}); err != nil {
		logging.LogError(err, "Failed to send message to grpc stream")
		closedConnection <- true
		return
	} else {
		for {
			//logging.LogDebug("in for loop waiting to receive from grpc translate from custom to mythic format")
			if input, err := stream.Recv(); err == io.EOF {
				logging.LogError(err, "got EOF from other side for handle custom format to mythic format")
				closedConnection <- true
				return
			} else if err != nil {
				logging.LogError(err, "Failed to read from stream for custom format to mythic format")
				closedConnection <- true
				return
			} else {
				if translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().TranslateCustomToMythicFormat != nil {
					sendMsg := translationstructs.TrCustomMessageToMythicC2FormatMessage{
						TranslationContainerName: input.GetTranslationContainerName(),
						C2Name:                   input.GetC2Name(),
						Message:                  input.GetMessage(),
						UUID:                     input.GetUUID(),
						MythicEncrypts:           input.GetMythicEncrypts(),
					}
					grpcCryptoKeys := input.GetCryptoKeys()
					cryptoKeys := make([]translationstructs.CryptoKeys, len(grpcCryptoKeys))
					for i := 0; i < len(cryptoKeys); i++ {
						cryptoKeys[i].DecKey = &grpcCryptoKeys[i].DecKey
						cryptoKeys[i].EncKey = &grpcCryptoKeys[i].EncKey
						cryptoKeys[i].Value = grpcCryptoKeys[i].Value
					}
					sendMsg.CryptoKeys = cryptoKeys
					response := translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().TranslateCustomToMythicFormat(sendMsg)
					sendResp := services.TrCustomMessageToMythicC2FormatMessageResponse{
						Success:                  response.Success,
						Error:                    response.Error,
						TranslationContainerName: input.GetTranslationContainerName(),
					}
					if messageBytes, err := json.Marshal(response.Message); err != nil {
						logging.LogError(err, "Failed to convert interface to bytes")
						sendResp.Success = false
						sendResp.Error = err.Error()
					} else {
						sendResp.Message = messageBytes
					}
					if err := stream.Send(&sendResp); err != nil {
						logging.LogError(err, "Failed to send response back to Mythic over grpc")
					}
				}
			}
		}
	}
}
func handleMythicToCustomFormat(wg *sync.WaitGroup, translationContainerName string, client services.TranslationContainerClient, closedConnection chan bool) {
	defer wg.Done()
	if stream, err := client.TranslateFromMythicToCustomFormat(context.Background()); err != nil {
		logging.LogError(err, "Failed to connect to grpc for TranslateFromCustomToMythicFormat")
		closedConnection <- true
		return
	} else if err := stream.Send(&services.TrMythicC2ToCustomMessageFormatMessageResponse{
		TranslationContainerName: translationContainerName,
	}); err != nil {
		logging.LogError(err, "Failed to send message to grpc stream")
		closedConnection <- true
		return
	} else {
		for {
			//logging.LogDebug("in for loop waiting to receive from grpc translate from custom to mythic format")
			if input, err := stream.Recv(); err == io.EOF {
				logging.LogError(err, "got EOF from other side for handle custom format to mythic format")
				closedConnection <- true
				return
			} else if err != nil {
				logging.LogError(err, "Failed to read from stream for custom format to mythic format")
				closedConnection <- true
				return
			} else {
				if translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().TranslateMythicToCustomFormat != nil {
					sendMsg := translationstructs.TrMythicC2ToCustomMessageFormatMessage{
						TranslationContainerName: input.GetTranslationContainerName(),
						C2Name:                   input.GetC2Name(),
						UUID:                     input.GetUUID(),
						MythicEncrypts:           input.GetMythicEncrypts(),
					}
					messageMap := map[string]interface{}{}
					if err = json.Unmarshal(input.GetMessage(), &messageMap); err != nil {
						logging.LogError(err, "Failed to unmarshal bytes into map")
						sendMsg.Message = messageMap
					} else {
						sendMsg.Message = messageMap
					}
					grpcCryptoKeys := input.GetCryptoKeys()
					cryptoKeys := make([]translationstructs.CryptoKeys, len(grpcCryptoKeys))
					for i := 0; i < len(cryptoKeys); i++ {
						cryptoKeys[i].DecKey = &grpcCryptoKeys[i].DecKey
						cryptoKeys[i].EncKey = &grpcCryptoKeys[i].EncKey
						cryptoKeys[i].Value = grpcCryptoKeys[i].Value
					}
					sendMsg.CryptoKeys = cryptoKeys
					response := translationstructs.AllTranslationData.Get(input.TranslationContainerName).GetPayloadDefinition().TranslateMythicToCustomFormat(sendMsg)
					sendResp := services.TrMythicC2ToCustomMessageFormatMessageResponse{
						Success:                  response.Success,
						Error:                    response.Error,
						TranslationContainerName: input.GetTranslationContainerName(),
						Message:                  response.Message,
					}
					if err := stream.Send(&sendResp); err != nil {
						logging.LogError(err, "Failed to send response back to Mythic over grpc")
					}
				}
			}
		}
	}
}

@echo off
echo Listing all .proto files...

set SCRIPT_DIR=%~dp0
echo Script is running from: %SCRIPT_DIR%

for /R proto %%F in (*.proto) do (
    echo Processing: %%F
	protoc -I=proto --proto_path=%SCRIPT_DIR% --go_out=interfaces --go-grpc_out=interfaces --grpc-gateway_out=interfaces --openapiv2_out=interfaces %%F
)

echo Done.
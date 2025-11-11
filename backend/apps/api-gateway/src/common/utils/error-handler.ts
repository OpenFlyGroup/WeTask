import { HttpException, HttpStatus } from '@nestjs/common';
import { RPCResponse } from '@libs/common';

export function handleRPCError(result: RPCResponse): never {
  const statusCode = result.statusCode || HttpStatus.INTERNAL_SERVER_ERROR;
  const errorMessage = result.error || 'Internal server error';
  throw new HttpException(errorMessage, statusCode);
}


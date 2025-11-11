import {
  WebSocketGateway,
  WebSocketServer,
  SubscribeMessage,
  OnGatewayConnection,
  OnGatewayDisconnect,
  MessageBody,
  ConnectedSocket,
} from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';
import { Injectable, Inject } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { RabbitMQEvents } from '@libs/common';

@WebSocketGateway({
  cors: {
    origin: '*',
  },
  namespace: '/',
})
@Injectable()
export class AppWebSocketGateway
  implements OnGatewayConnection, OnGatewayDisconnect
{
  @WebSocketServer()
  server: Server;

  private userRooms: Map<number, Set<string>> = new Map();

  constructor(
    @Inject('TASKS_SERVICE') private tasksClient: ClientProxy,
    @Inject('BOARDS_SERVICE') private boardsClient: ClientProxy,
    @Inject('TEAMS_SERVICE') private teamsClient: ClientProxy,
  ) {}

  handleConnection(client: Socket) {
    console.log(`Client connected: ${client.id}`);
  }

  handleDisconnect(client: Socket) {
    console.log(`Client disconnected: ${client.id}`);
    // Clean up user rooms
    this.userRooms.forEach((rooms, userId) => {
      rooms.delete(client.id);
      if (rooms.size === 0) {
        this.userRooms.delete(userId);
      }
    });
  }

  @SubscribeMessage('join:board')
  handleJoinBoard(
    @MessageBody() data: { boardId: number; userId: number },
    @ConnectedSocket() client: Socket,
  ) {
    const room = `board:${data.boardId}`;
    void client.join(room);
    console.log(`Client ${client.id} joined board ${data.boardId}`);
  }

  @SubscribeMessage('leave:board')
  handleLeaveBoard(
    @MessageBody() data: { boardId: number },
    @ConnectedSocket() client: Socket,
  ) {
    const room = `board:${data.boardId}`;
    void client.leave(room);
    console.log(`Client ${client.id} left board ${data.boardId}`);
  }

  @SubscribeMessage('join:team')
  handleJoinTeam(
    @MessageBody() data: { teamId: number; userId: number },
    @ConnectedSocket() client: Socket,
  ) {
    const room = `team:${data.teamId}`;
    void client.join(room);

    if (!this.userRooms.has(data.userId)) {
      this.userRooms.set(data.userId, new Set());
    }
    this.userRooms.get(data.userId)!.add(client.id);

    console.log(`Client ${client.id} joined team ${data.teamId}`);
  }

  // Methods to emit events to clients
  emitTaskCreated(boardId: number, task: any) {
    this.server.to(`board:${boardId}`).emit(RabbitMQEvents.TASK_CREATED, task);
  }

  emitTaskUpdated(boardId: number, task: any) {
    this.server.to(`board:${boardId}`).emit(RabbitMQEvents.TASK_UPDATED, task);
  }

  emitTaskDeleted(boardId: number, taskId: number) {
    this.server
      .to(`board:${boardId}`)
      .emit(RabbitMQEvents.TASK_DELETED, { taskId });
  }

  emitBoardUpdated(teamId: number, board: any) {
    this.server.to(`team:${teamId}`).emit(RabbitMQEvents.BOARD_UPDATED, board);
  }

  emitTeamMemberAdded(teamId: number, member: any) {
    this.server
      .to(`team:${teamId}`)
      .emit(RabbitMQEvents.TEAM_MEMBER_ADDED, member);
  }

  emitTeamMemberRemoved(teamId: number, userId: number) {
    this.server
      .to(`team:${teamId}`)
      .emit(RabbitMQEvents.TEAM_MEMBER_REMOVED, { userId });
  }
}

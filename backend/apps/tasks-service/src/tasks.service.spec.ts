import { Test, TestingModule } from '@nestjs/testing';
import { getModelToken } from '@nestjs/mongoose';
import { TasksService } from './tasks.service';
import { PrismaService } from '@libs/prisma';
import { Comment } from './schemas/comment.schema';
import { ActivityLog } from './schemas/activity-log.schema';

describe('TasksService', () => {
  let service: TasksService;
  let prismaService: jest.Mocked<PrismaService>;
  let commentModel: jest.Mocked<any>;
  let activityLogModel: jest.Mocked<any>;

  const mockTask = {
    id: 1,
    title: 'Test Task',
    description: 'Test Description',
    columnId: 1,
    assignedTo: 1,
    priority: 'medium',
    status: 'todo',
    createdAt: new Date(),
    updatedAt: new Date(),
    user: {
      id: 1,
      email: 'test@example.com',
      name: 'Test User',
    },
    column: {
      id: 1,
      title: 'Test Column',
      board: {
        id: 1,
        title: 'Test Board',
      },
    },
  };

  const mockComment = {
    _id: 'comment-id',
    taskId: 1,
    userId: 1,
    message: 'Test comment',
    createdAt: new Date(),
    updatedAt: new Date(),
  };

  const mockActivityLog = {
    _id: 'log-id',
    entityType: 'task',
    entityId: 1,
    action: 'created',
    userId: 1,
    metadata: {},
    timestamp: new Date(),
  };

  beforeEach(async () => {
    const mockPrismaService = {
      task: {
        create: jest.fn(),
        findUnique: jest.fn(),
        update: jest.fn(),
        delete: jest.fn(),
      },
      board: {
        findUnique: jest.fn(),
      },
    };

    const mockCommentModel = {
      create: jest.fn(),
      find: jest.fn().mockReturnValue({
        sort: jest.fn().mockReturnValue({
          exec: jest.fn(),
        }),
      }),
    };

    const mockActivityLogModel = {
      create: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        TasksService,
        {
          provide: PrismaService,
          useValue: mockPrismaService,
        },
        {
          provide: getModelToken(Comment.name),
          useValue: mockCommentModel,
        },
        {
          provide: getModelToken(ActivityLog.name),
          useValue: mockActivityLogModel,
        },
      ],
    }).compile();

    service = module.get<TasksService>(TasksService);
    prismaService = module.get(PrismaService);
    commentModel = module.get(getModelToken(Comment.name));
    activityLogModel = module.get(getModelToken(ActivityLog.name));
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('create', () => {
    it('should create a task successfully', async () => {
      prismaService.task.create.mockResolvedValue(mockTask);
      activityLogModel.create.mockResolvedValue(mockActivityLog);

      const result = await service.create({
        title: 'Test Task',
        description: 'Test Description',
        columnId: 1,
        assignedTo: 1,
        priority: 'high',
      });

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockTask);
      expect(prismaService.task.create).toHaveBeenCalledWith({
        data: {
          title: 'Test Task',
          description: 'Test Description',
          columnId: 1,
          assignedTo: 1,
          priority: 'high',
          status: 'todo',
        },
        include: {
          user: {
            select: {
              id: true,
              email: true,
              name: true,
            },
          },
          column: {
            include: {
              board: true,
            },
          },
        },
      });
      expect(activityLogModel.create).toHaveBeenCalled();
    });

    it('should use default priority if not provided', async () => {
      prismaService.task.create.mockResolvedValue(mockTask);
      activityLogModel.create.mockResolvedValue(mockActivityLog);

      await service.create({
        title: 'Test Task',
        columnId: 1,
      });

      expect(prismaService.task.create).toHaveBeenCalledWith(
        expect.objectContaining({
          data: expect.objectContaining({
            priority: 'medium',
          }),
        }),
      );
    });
  });

  describe('getById', () => {
    it('should return task by id', async () => {
      prismaService.task.findUnique.mockResolvedValue(mockTask);

      const result = await service.getById(1);

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockTask);
    });

    it('should return error if task not found', async () => {
      prismaService.task.findUnique.mockResolvedValue(null);

      const result = await service.getById(999);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Task not found');
      expect(result.statusCode).toBe(404);
    });
  });

  describe('getByBoard', () => {
    it('should return all tasks for a board', async () => {
      const mockBoard = {
        id: 1,
        title: 'Test Board',
        columns: [
          {
            id: 1,
            tasks: [mockTask],
          },
        ],
      };
      prismaService.board.findUnique.mockResolvedValue(mockBoard);

      const result = await service.getByBoard(1);

      expect(result.success).toBe(true);
      expect(result.data).toEqual([mockTask]);
    });

    it('should return error if board not found', async () => {
      prismaService.board.findUnique.mockResolvedValue(null);

      const result = await service.getByBoard(999);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Board not found');
      expect(result.statusCode).toBe(404);
    });
  });

  describe('update', () => {
    it('should update task successfully', async () => {
      const updatedTask = { ...mockTask, title: 'Updated Task' };
      prismaService.task.update.mockResolvedValue(updatedTask);
      activityLogModel.create.mockResolvedValue(mockActivityLog);

      const result = await service.update(1, {
        title: 'Updated Task',
        description: 'Updated Description',
      });

      expect(result.success).toBe(true);
      expect(result.data.title).toBe('Updated Task');
      expect(activityLogModel.create).toHaveBeenCalled();
    });

    it('should handle undefined description', async () => {
      const updatedTask = { ...mockTask, description: null };
      prismaService.task.update.mockResolvedValue(updatedTask);
      activityLogModel.create.mockResolvedValue(mockActivityLog);

      await service.update(1, {
        description: undefined,
      });

      // When description is undefined, it should not be included in updateData
      const updateCall = prismaService.task.update.mock.calls[0][0];
      expect(updateCall.data).not.toHaveProperty('description');
    });
  });

  describe('delete', () => {
    it('should delete task successfully', async () => {
      prismaService.task.findUnique.mockResolvedValue(mockTask);
      prismaService.task.delete.mockResolvedValue(mockTask);
      activityLogModel.create.mockResolvedValue(mockActivityLog);

      const result = await service.delete(1);

      expect(result.success).toBe(true);
      expect(result.data.message).toBe('Task deleted successfully');
      expect(activityLogModel.create).toHaveBeenCalled();
    });

    it('should return error if task not found', async () => {
      prismaService.task.findUnique.mockResolvedValue(null);

      const result = await service.delete(999);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Task not found');
      expect(result.statusCode).toBe(404);
    });
  });

  describe('move', () => {
    it('should move task to another column', async () => {
      const movedTask = { ...mockTask, columnId: 2 };
      prismaService.task.findUnique.mockResolvedValue(mockTask);
      prismaService.task.update.mockResolvedValue(movedTask);
      activityLogModel.create.mockResolvedValue(mockActivityLog);

      const result = await service.move(1, 2);

      expect(result.success).toBe(true);
      expect(result.data.columnId).toBe(2);
      expect(activityLogModel.create).toHaveBeenCalledWith(
        expect.objectContaining({
          action: 'moved',
          metadata: {
            fromColumnId: 1,
            toColumnId: 2,
          },
        }),
      );
    });

    it('should return error if task not found', async () => {
      prismaService.task.findUnique.mockResolvedValue(null);

      const result = await service.move(999, 2);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Task not found');
      expect(result.statusCode).toBe(404);
    });
  });

  describe('addComment', () => {
    it('should add comment to task', async () => {
      commentModel.create.mockResolvedValue(mockComment);
      activityLogModel.create.mockResolvedValue(mockActivityLog);

      const result = await service.addComment({
        taskId: 1,
        userId: 1,
        message: 'Test comment',
      });

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockComment);
      expect(commentModel.create).toHaveBeenCalledWith({
        taskId: 1,
        userId: 1,
        message: 'Test comment',
      });
      expect(activityLogModel.create).toHaveBeenCalled();
    });
  });

  describe('getComments', () => {
    it('should return all comments for a task', async () => {
      const execMock = jest.fn().mockResolvedValue([mockComment]);
      commentModel.find.mockReturnValue({
        sort: jest.fn().mockReturnValue({
          exec: execMock,
        }),
      });

      const result = await service.getComments(1);

      expect(result.success).toBe(true);
      expect(result.data).toEqual([mockComment]);
      expect(commentModel.find).toHaveBeenCalledWith({ taskId: 1 });
    });
  });
});


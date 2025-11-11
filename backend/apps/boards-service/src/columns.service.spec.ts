import { Test, TestingModule } from '@nestjs/testing';
import { ColumnsService } from './columns.service';
import { PrismaService } from '@libs/prisma';

describe('ColumnsService', () => {
  let service: ColumnsService;
  let prismaService: jest.Mocked<PrismaService>;

  const mockColumn = {
    id: 1,
    title: 'Test Column',
    boardId: 1,
    order: 0,
    createdAt: new Date(),
    updatedAt: new Date(),
  };

  const mockColumnWithTasks = {
    ...mockColumn,
    tasks: [],
  };

  beforeEach(async () => {
    const mockPrismaService = {
      column: {
        create: jest.fn(),
        findMany: jest.fn(),
        update: jest.fn(),
        delete: jest.fn(),
      },
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        ColumnsService,
        {
          provide: PrismaService,
          useValue: mockPrismaService,
        },
      ],
    }).compile();

    service = module.get<ColumnsService>(ColumnsService);
    prismaService = module.get(PrismaService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('create', () => {
    it('should create a column successfully', async () => {
      prismaService.column.create.mockResolvedValue(mockColumn);

      const result = await service.create({
        title: 'Test Column',
        boardId: 1,
        order: 0,
      });

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockColumn);
      expect(prismaService.column.create).toHaveBeenCalledWith({
        data: {
          title: 'Test Column',
          boardId: 1,
          order: 0,
        },
      });
    });
  });

  describe('getByBoard', () => {
    it('should return all columns for a board', async () => {
      prismaService.column.findMany.mockResolvedValue([mockColumnWithTasks]);

      const result = await service.getByBoard(1);

      expect(result.success).toBe(true);
      expect(result.data).toEqual([mockColumnWithTasks]);
      expect(prismaService.column.findMany).toHaveBeenCalledWith({
        where: { boardId: 1 },
        orderBy: {
          order: 'asc',
        },
        include: {
          tasks: {
            include: {
              user: {
                select: {
                  id: true,
                  email: true,
                  name: true,
                },
              },
            },
          },
        },
      });
    });
  });

  describe('update', () => {
    it('should update column title', async () => {
      const updatedColumn = { ...mockColumn, title: 'Updated Title' };
      prismaService.column.update.mockResolvedValue(updatedColumn);

      const result = await service.update(1, { title: 'Updated Title' });

      expect(result.success).toBe(true);
      expect(result.data.title).toBe('Updated Title');
      expect(prismaService.column.update).toHaveBeenCalledWith({
        where: { id: 1 },
        data: { title: 'Updated Title' },
      });
    });

    it('should update column order', async () => {
      const updatedColumn = { ...mockColumn, order: 1 };
      prismaService.column.update.mockResolvedValue(updatedColumn);

      const result = await service.update(1, { order: 1 });

      expect(result.success).toBe(true);
      expect(result.data.order).toBe(1);
      expect(prismaService.column.update).toHaveBeenCalledWith({
        where: { id: 1 },
        data: { order: 1 },
      });
    });

    it('should update both title and order', async () => {
      const updatedColumn = { ...mockColumn, title: 'Updated Title', order: 2 };
      prismaService.column.update.mockResolvedValue(updatedColumn);

      const result = await service.update(1, { title: 'Updated Title', order: 2 });

      expect(result.success).toBe(true);
      expect(result.data.title).toBe('Updated Title');
      expect(result.data.order).toBe(2);
    });
  });

  describe('delete', () => {
    it('should delete column successfully', async () => {
      prismaService.column.delete.mockResolvedValue(mockColumn);

      const result = await service.delete(1);

      expect(result.success).toBe(true);
      expect(result.data.message).toBe('Column deleted successfully');
      expect(prismaService.column.delete).toHaveBeenCalledWith({
        where: { id: 1 },
      });
    });
  });
});


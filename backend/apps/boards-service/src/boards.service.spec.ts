import { Test, TestingModule } from '@nestjs/testing';
import { BoardsService } from './boards.service';
import { PrismaService } from '@libs/prisma';

describe('BoardsService', () => {
  let service: BoardsService;
  let prismaService: jest.Mocked<PrismaService>;

  const mockBoard = {
    id: 1,
    title: 'Test Board',
    teamId: 1,
    createdAt: new Date(),
    updatedAt: new Date(),
    team: {
      id: 1,
      name: 'Test Team',
    },
    columns: [],
  };

  const mockTeam = {
    id: 1,
    name: 'Test Team',
    boards: [mockBoard],
    members: [],
  };

  beforeEach(async () => {
    const mockPrismaService = {
      board: {
        create: jest.fn(),
        findUnique: jest.fn(),
        findMany: jest.fn(),
        update: jest.fn(),
        delete: jest.fn(),
      },
      team: {
        findMany: jest.fn(),
      },
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        BoardsService,
        {
          provide: PrismaService,
          useValue: mockPrismaService,
        },
      ],
    }).compile();

    service = module.get<BoardsService>(BoardsService);
    prismaService = module.get(PrismaService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('create', () => {
    it('should create a board successfully', async () => {
      prismaService.board.create.mockResolvedValue(mockBoard);

      const result = await service.create({
        title: 'Test Board',
        teamId: 1,
      });

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockBoard);
      expect(prismaService.board.create).toHaveBeenCalledWith({
        data: {
          title: 'Test Board',
          teamId: 1,
        },
        include: {
          team: true,
          columns: {
            orderBy: {
              order: 'asc',
            },
          },
        },
      });
    });
  });

  describe('getAllByUser', () => {
    it('should return all boards for a user', async () => {
      prismaService.team.findMany.mockResolvedValue([mockTeam]);

      const result = await service.getAllByUser(1);

      expect(result.success).toBe(true);
      expect(result.data).toEqual([mockBoard]);
      expect(prismaService.team.findMany).toHaveBeenCalledWith({
        where: {
          members: {
            some: {
              userId: 1,
            },
          },
        },
        include: {
          boards: {
            include: {
              columns: {
                orderBy: {
                  order: 'asc',
                },
              },
            },
          },
        },
      });
    });
  });

  describe('getById', () => {
    it('should return board by id', async () => {
      const boardWithDetails = {
        ...mockBoard,
        team: {
          ...mockBoard.team,
          members: [],
        },
        columns: [],
      };
      prismaService.board.findUnique.mockResolvedValue(boardWithDetails);

      const result = await service.getById(1);

      expect(result.success).toBe(true);
      expect(result.data).toEqual(boardWithDetails);
    });

    it('should return error if board not found', async () => {
      prismaService.board.findUnique.mockResolvedValue(null);

      const result = await service.getById(999);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Board not found');
      expect(result.statusCode).toBe(404);
    });
  });

  describe('getByTeam', () => {
    it('should return all boards for a team', async () => {
      prismaService.board.findMany.mockResolvedValue([mockBoard]);

      const result = await service.getByTeam(1);

      expect(result.success).toBe(true);
      expect(result.data).toEqual([mockBoard]);
      expect(prismaService.board.findMany).toHaveBeenCalledWith({
        where: { teamId: 1 },
        include: {
          columns: {
            orderBy: {
              order: 'asc',
            },
          },
        },
      });
    });
  });

  describe('update', () => {
    it('should update board title', async () => {
      const updatedBoard = { ...mockBoard, title: 'Updated Title' };
      prismaService.board.update.mockResolvedValue(updatedBoard);

      const result = await service.update(1, { title: 'Updated Title' });

      expect(result.success).toBe(true);
      expect(result.data.title).toBe('Updated Title');
      expect(prismaService.board.update).toHaveBeenCalledWith({
        where: { id: 1 },
        data: { title: 'Updated Title' },
        include: {
          columns: {
            orderBy: {
              order: 'asc',
            },
          },
        },
      });
    });
  });

  describe('delete', () => {
    it('should delete board successfully', async () => {
      prismaService.board.delete.mockResolvedValue(mockBoard);

      const result = await service.delete(1);

      expect(result.success).toBe(true);
      expect(result.data.message).toBe('Board deleted successfully');
      expect(prismaService.board.delete).toHaveBeenCalledWith({
        where: { id: 1 },
      });
    });
  });
});


import { Test, TestingModule } from '@nestjs/testing';
import { TeamsService } from './teams.service';
import { PrismaService } from '@libs/prisma';

describe('TeamsService', () => {
  let service: TeamsService;
  let prismaService: jest.Mocked<PrismaService>;

  const mockTeam = {
    id: 1,
    name: 'Test Team',
    createdAt: new Date(),
    updatedAt: new Date(),
    members: [
      {
        id: 1,
        userId: 1,
        teamId: 1,
        role: 'owner',
        user: {
          id: 1,
          email: 'test@example.com',
          name: 'Test User',
        },
      },
    ],
  };

  const mockTeamMember = {
    id: 1,
    userId: 1,
    teamId: 1,
    role: 'member',
    user: {
      id: 1,
      email: 'test@example.com',
      name: 'Test User',
    },
  };

  beforeEach(async () => {
    const mockPrismaService = {
      team: {
        create: jest.fn(),
        findUnique: jest.fn(),
        findMany: jest.fn(),
      },
      teamMember: {
        findUnique: jest.fn(),
        create: jest.fn(),
        delete: jest.fn(),
      },
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        TeamsService,
        {
          provide: PrismaService,
          useValue: mockPrismaService,
        },
      ],
    }).compile();

    service = module.get<TeamsService>(TeamsService);
    prismaService = module.get(PrismaService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('create', () => {
    it('should create a team successfully', async () => {
      prismaService.team.create.mockResolvedValue(mockTeam);

      const result = await service.create({
        name: 'Test Team',
        ownerId: 1,
      });

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockTeam);
      expect(prismaService.team.create).toHaveBeenCalledWith({
        data: {
          name: 'Test Team',
          members: {
            create: {
              userId: 1,
              role: 'owner',
            },
          },
        },
        include: {
          members: {
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

  describe('getById', () => {
    it('should return team by id', async () => {
      prismaService.team.findUnique.mockResolvedValue(mockTeam);

      const result = await service.getById(1);

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockTeam);
    });

    it('should return error if team not found', async () => {
      prismaService.team.findUnique.mockResolvedValue(null);

      const result = await service.getById(999);

      expect(result.success).toBe(false);
      expect(result.error).toBe('Team not found');
      expect(result.statusCode).toBe(404);
    });
  });

  describe('getUserTeams', () => {
    it('should return all teams for a user', async () => {
      prismaService.team.findMany.mockResolvedValue([mockTeam]);

      const result = await service.getUserTeams(1);

      expect(result.success).toBe(true);
      expect(result.data).toEqual([mockTeam]);
      expect(prismaService.team.findMany).toHaveBeenCalledWith({
        where: {
          members: {
            some: {
              userId: 1,
            },
          },
        },
        include: {
          members: {
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

  describe('addMember', () => {
    it('should add member to team successfully', async () => {
      prismaService.team.findUnique.mockResolvedValue(mockTeam);
      prismaService.teamMember.findUnique.mockResolvedValue(null);
      prismaService.teamMember.create.mockResolvedValue(mockTeamMember);

      const result = await service.addMember({
        teamId: 1,
        userId: 2,
        role: 'member',
      });

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockTeamMember);
    });

    it('should return error if team not found', async () => {
      prismaService.team.findUnique.mockResolvedValue(null);

      const result = await service.addMember({
        teamId: 999,
        userId: 2,
      });

      expect(result.success).toBe(false);
      expect(result.error).toBe('Team not found');
      expect(result.statusCode).toBe(404);
    });

    it('should return error if user is already a member', async () => {
      prismaService.team.findUnique.mockResolvedValue(mockTeam);
      prismaService.teamMember.findUnique.mockResolvedValue(mockTeamMember);

      const result = await service.addMember({
        teamId: 1,
        userId: 1,
      });

      expect(result.success).toBe(false);
      expect(result.error).toBe('User is already a member of this team');
      expect(result.statusCode).toBe(409);
    });

    it('should use default role if not provided', async () => {
      prismaService.team.findUnique.mockResolvedValue(mockTeam);
      prismaService.teamMember.findUnique.mockResolvedValue(null);
      prismaService.teamMember.create.mockResolvedValue(mockTeamMember);

      await service.addMember({
        teamId: 1,
        userId: 2,
      });

      expect(prismaService.teamMember.create).toHaveBeenCalledWith({
        data: {
          teamId: 1,
          userId: 2,
          role: 'member',
        },
        include: {
          user: {
            select: {
              id: true,
              email: true,
              name: true,
            },
          },
        },
      });
    });
  });

  describe('removeMember', () => {
    it('should remove member from team successfully', async () => {
      prismaService.teamMember.findUnique.mockResolvedValue(mockTeamMember);
      prismaService.teamMember.delete.mockResolvedValue(mockTeamMember);

      const result = await service.removeMember({
        teamId: 1,
        userId: 1,
      });

      expect(result.success).toBe(true);
      expect(result.data.message).toBe('Member removed successfully');
      expect(prismaService.teamMember.delete).toHaveBeenCalled();
    });

    it('should return error if member not found', async () => {
      prismaService.teamMember.findUnique.mockResolvedValue(null);

      const result = await service.removeMember({
        teamId: 1,
        userId: 999,
      });

      expect(result.success).toBe(false);
      expect(result.error).toBe('Member not found');
      expect(result.statusCode).toBe(404);
    });
  });
});


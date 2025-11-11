import { Test, TestingModule } from '@nestjs/testing';
import { UsersService } from './users.service';
import { PrismaService } from '@libs/prisma';

describe('UsersService', () => {
  let service: UsersService;
  let prismaService: jest.Mocked<PrismaService>;

  const mockUser = {
    id: 1,
    email: 'test@example.com',
    name: 'Test User',
    createdAt: new Date(),
    updatedAt: new Date(),
  };

  beforeEach(async () => {
    const mockPrismaService = {
      user: {
        findUnique: jest.fn(),
        update: jest.fn(),
      },
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        UsersService,
        {
          provide: PrismaService,
          useValue: mockPrismaService,
        },
      ],
    }).compile();

    service = module.get<UsersService>(UsersService);
    prismaService = module.get(PrismaService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('getById', () => {
    it('should return user by id', async () => {
      prismaService.user.findUnique.mockResolvedValue(mockUser);

      const result = await service.getById(1);

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockUser);
      expect(prismaService.user.findUnique).toHaveBeenCalledWith({
        where: { id: 1 },
        select: {
          id: true,
          email: true,
          name: true,
          createdAt: true,
          updatedAt: true,
        },
      });
    });

    it('should return error if user not found', async () => {
      prismaService.user.findUnique.mockResolvedValue(null);

      const result = await service.getById(999);

      expect(result.success).toBe(false);
      expect(result.error).toBe('User not found');
      expect(result.statusCode).toBe(404);
    });
  });

  describe('getByEmail', () => {
    it('should return user by email', async () => {
      prismaService.user.findUnique.mockResolvedValue(mockUser);

      const result = await service.getByEmail('test@example.com');

      expect(result.success).toBe(true);
      expect(result.data).toEqual(mockUser);
      expect(prismaService.user.findUnique).toHaveBeenCalledWith({
        where: { email: 'test@example.com' },
        select: {
          id: true,
          email: true,
          name: true,
          createdAt: true,
          updatedAt: true,
        },
      });
    });

    it('should return error if user not found', async () => {
      prismaService.user.findUnique.mockResolvedValue(null);

      const result = await service.getByEmail('notfound@example.com');

      expect(result.success).toBe(false);
      expect(result.error).toBe('User not found');
      expect(result.statusCode).toBe(404);
    });
  });

  describe('update', () => {
    it('should update user name', async () => {
      const updatedUser = { ...mockUser, name: 'Updated Name' };
      prismaService.user.update.mockResolvedValue(updatedUser);

      const result = await service.update(1, { name: 'Updated Name' });

      expect(result.success).toBe(true);
      expect(result.data.name).toBe('Updated Name');
      expect(prismaService.user.update).toHaveBeenCalledWith({
        where: { id: 1 },
        data: { name: 'Updated Name' },
        select: {
          id: true,
          email: true,
          name: true,
          createdAt: true,
          updatedAt: true,
        },
      });
    });

    it('should update user email', async () => {
      const updatedUser = { ...mockUser, email: 'newemail@example.com' };
      prismaService.user.update.mockResolvedValue(updatedUser);

      const result = await service.update(1, { email: 'newemail@example.com' });

      expect(result.success).toBe(true);
      expect(result.data.email).toBe('newemail@example.com');
    });

    it('should update both name and email', async () => {
      const updatedUser = {
        ...mockUser,
        name: 'Updated Name',
        email: 'newemail@example.com',
      };
      prismaService.user.update.mockResolvedValue(updatedUser);

      const result = await service.update(1, {
        name: 'Updated Name',
        email: 'newemail@example.com',
      });

      expect(result.success).toBe(true);
      expect(result.data.name).toBe('Updated Name');
      expect(result.data.email).toBe('newemail@example.com');
    });
  });
});


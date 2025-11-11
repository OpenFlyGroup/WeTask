import { Test, TestingModule } from '@nestjs/testing';
import { JwtService } from '@nestjs/jwt';
import { AuthService } from './auth.service';
import { DatabaseService } from '@libs/database';

jest.mock('bcryptjs', () => ({
  hash: jest.fn(),
  compare: jest.fn(),
}));

jest.mock('crypto', () => ({
  randomBytes: jest.fn(() => ({
    toString: jest.fn(() => 'mock-refresh-token'),
  })),
}));

import * as bcrypt from 'bcryptjs';

describe('AuthService', () => {
  let service: AuthService;
  let databaseService: jest.Mocked<DatabaseService>;
  let jwtService: jest.Mocked<JwtService>;

  const mockUser = {
    id: 1,
    email: 'test@example.com',
    password: 'hashed-password',
    name: 'Test User',
    createdAt: new Date(),
  };

  const mockRefreshToken = {
    id: 1,
    token: 'refresh-token',
    userId: 1,
    expiresAt: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000),
    user: mockUser,
  };

  beforeEach(async () => {
    const mockDatabaseService = {
      user: {
        findOne: jest.fn(),
        save: jest.fn(),
      },
      refreshToken: {
        findOne: jest.fn(),
        save: jest.fn(),
        remove: jest.fn(),
      },
    };

    const mockJwtService = {
      sign: jest.fn(() => 'access-token'),
      verify: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        AuthService,
        {
          provide: DatabaseService,
          useValue: mockDatabaseService,
        },
        {
          provide: JwtService,
          useValue: mockJwtService,
        },
      ],
    }).compile();

    service = module.get<AuthService>(AuthService);
    databaseService = module.get(DatabaseService);
    jwtService = module.get(JwtService);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('register', () => {
    it('should register a new user successfully', async () => {
      databaseService.user.findOne.mockResolvedValue(null);
      (bcrypt.hash as jest.Mock).mockResolvedValue('hashed-password');
      databaseService.user.save.mockResolvedValue({
        ...mockUser,
        password: 'hashed-password',
      });
      databaseService.refreshToken.save.mockResolvedValue(mockRefreshToken);

      const result = await service.register({
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User',
      });

      expect(result.success).toBe(true);
      expect(result.data).toHaveProperty('user');
      expect(result.data).toHaveProperty('accessToken');
      expect(result.data).toHaveProperty('refreshToken');
      expect(databaseService.user.findOne).toHaveBeenCalledWith({
        where: { email: 'test@example.com' },
      });
      expect(bcrypt.hash).toHaveBeenCalledWith('password123', 10);
    });

    it('should return error if user already exists', async () => {
      databaseService.user.findOne.mockResolvedValue(mockUser);

      const result = await service.register({
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User',
      });

      expect(result.success).toBe(false);
      expect(result.error).toBe('User with this email already exists');
      expect(result.statusCode).toBe(409);
    });
  });

  describe('login', () => {
    it('should login user successfully', async () => {
      databaseService.user.findOne.mockResolvedValue(mockUser);
      (bcrypt.compare as jest.Mock).mockResolvedValue(true);
      databaseService.refreshToken.save.mockResolvedValue(mockRefreshToken);

      const result = await service.login({
        email: 'test@example.com',
        password: 'password123',
      });

      expect(result.success).toBe(true);
      expect(result.data).toHaveProperty('user');
      expect(result.data).toHaveProperty('accessToken');
      expect(result.data).toHaveProperty('refreshToken');
      expect(bcrypt.compare).toHaveBeenCalledWith('password123', 'hashed-password');
    });

    it('should return error if user not found', async () => {
      databaseService.user.findOne.mockResolvedValue(null);

      const result = await service.login({
        email: 'test@example.com',
        password: 'password123',
      });

      expect(result.success).toBe(false);
      expect(result.error).toBe('Invalid credentials');
      expect(result.statusCode).toBe(401);
    });

    it('should return error if password is invalid', async () => {
      databaseService.user.findOne.mockResolvedValue(mockUser);
      (bcrypt.compare as jest.Mock).mockResolvedValue(false);

      const result = await service.login({
        email: 'test@example.com',
        password: 'wrong-password',
      });

      expect(result.success).toBe(false);
      expect(result.error).toBe('Invalid credentials');
      expect(result.statusCode).toBe(401);
    });
  });

  describe('refresh', () => {
    it('should refresh tokens successfully', async () => {
      databaseService.refreshToken.findOne.mockResolvedValue(mockRefreshToken);
      databaseService.refreshToken.remove.mockResolvedValue(mockRefreshToken);
      databaseService.refreshToken.save.mockResolvedValue(mockRefreshToken);

      const result = await service.refresh('refresh-token');

      expect(result.success).toBe(true);
      expect(result.data).toHaveProperty('accessToken');
      expect(result.data).toHaveProperty('refreshToken');
      expect(databaseService.refreshToken.remove).toHaveBeenCalled();
    });

    it('should return error if refresh token not found', async () => {
      databaseService.refreshToken.findOne.mockResolvedValue(null);

      const result = await service.refresh('invalid-token');

      expect(result.success).toBe(false);
      expect(result.error).toBe('Invalid or expired refresh token');
      expect(result.statusCode).toBe(401);
    });

    it('should return error if refresh token expired', async () => {
      const expiredToken = {
        ...mockRefreshToken,
        expiresAt: new Date(Date.now() - 1000),
      };
      databaseService.refreshToken.findOne.mockResolvedValue(expiredToken);

      const result = await service.refresh('expired-token');

      expect(result.success).toBe(false);
      expect(result.error).toBe('Invalid or expired refresh token');
      expect(result.statusCode).toBe(401);
    });
  });

  describe('validateToken', () => {
    it('should validate token successfully', async () => {
      jwtService.verify.mockReturnValue({ sub: 1 });
      databaseService.user.findOne.mockResolvedValue(mockUser);

      const result = await service.validateToken('valid-token');

      expect(result.success).toBe(true);
      expect(result.data).toHaveProperty('id');
      expect(result.data).toHaveProperty('email');
      expect(result.data).toHaveProperty('name');
    });

    it('should return error if token is invalid', async () => {
      jwtService.verify.mockImplementation(() => {
        throw new Error('Invalid token');
      });

      const result = await service.validateToken('invalid-token');

      expect(result.success).toBe(false);
      expect(result.error).toBe('Invalid token');
      expect(result.statusCode).toBe(401);
    });

    it('should return error if user not found', async () => {
      jwtService.verify.mockReturnValue({ sub: 999 });
      databaseService.user.findOne.mockResolvedValue(null);

      const result = await service.validateToken('valid-token');

      expect(result.success).toBe(false);
      expect(result.error).toBe('User not found');
      expect(result.statusCode).toBe(401);
    });
  });
});


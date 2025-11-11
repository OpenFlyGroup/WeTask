import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { Document } from 'mongoose';

@Schema({ timestamps: true })
export class ActivityLog extends Document {
  @Prop({ required: true })
  entityType: string; // 'task', 'board', 'team', etc.

  @Prop({ required: true })
  entityId: number;

  @Prop({ required: true })
  action: string; // 'created', 'updated', 'deleted', 'moved', etc.

  @Prop({ required: true })
  userId: number;

  @Prop({ type: Object })
  metadata?: Record<string, any>;

  @Prop()
  timestamp?: Date;
}

export const ActivityLogSchema = SchemaFactory.createForClass(ActivityLog);


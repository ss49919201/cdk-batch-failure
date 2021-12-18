import { Stack, StackProps, Duration } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as lambda from '@aws-cdk/aws-lambda-go-alpha';
import { SqsEventSource } from 'aws-cdk-lib/aws-lambda-event-sources';

export class BatchStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const deadLetterQueue = new sqs.Queue(this, "deadLetterQueue", {
      queueName: "deadLetterQueue",
    })

    // 一度でもLambda関数内で処理に失敗したメッセージはデッドレターキューに移動
    const queue = new sqs.Queue(this, "queue", {
      queueName: "queue",
      deadLetterQueue: {
        queue: deadLetterQueue,
        maxReceiveCount: 1,
      }
    })

    // イベントソースとなるQueue
    const source = new SqsEventSource(queue, {
      reportBatchItemFailures: true, // バッチ内のメッセージの一部を失敗したメッセージとしてキューに返す
    })
    new lambda.GoFunction(this, 'RandomResult', {
      entry: 'lambda/go',
      events: [source],
    })
  }
}

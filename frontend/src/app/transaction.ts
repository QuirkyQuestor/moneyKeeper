import { Timestamp } from "rxjs";

export interface Transaction {
  transactionId: string;
  accountFrom: string;
  date: Date;
  amount: number;
  accountTo: string;
  memo: string;
  categoryId: string;
  transferTransactionId: string;
}

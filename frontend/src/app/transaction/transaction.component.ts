import { Component, OnInit, Inject } from "@angular/core";
import { Account } from "../account";
import { AccountService } from "../account.service";
import { Category } from "../category";
import { CategoryService } from "../category.service";
import { Transaction } from "../transaction";
import { TransactionService } from "../transaction.service";
import {
  MatDialog,
  MAT_DIALOG_DATA,
  MatDialogRef,
} from "@angular/material/dialog";

@Component({
  selector: "app-transaction",
  templateUrl: "./transaction.component.html",
  styleUrls: ["./transaction.component.css"],
})
export class TransactionComponent implements OnInit {
  transactions: Transaction[] = [];
  categories: Category[] = [];
  accounts: Account[] = [];
  accountSelected: Account = {
    accountId: "",
    typeId: "",
    name: "",
    description: "",
    active: false,
  };

  constructor(
    private transactionService: TransactionService,
    private categoryService: CategoryService,
    private accountService: AccountService,
    public dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.getCategories();
    this.getAccounts();
    this.getTransactions();
  }

  getTransactions(): void {
    this.transactionService
      .getTransactions()
      .subscribe((transactions) => (this.transactions = transactions));
  }

  addTransaction(
    accountFrom: string,
    date: string,
    amount: string,
    accountTo: string,
    memo: string,
    categoryId: string,
    transferTransactionId: string
  ): void {
    // name = name.trim();
    // if (!name) {
    //   return;
    // }
    this.transactionService
      .addTransaction({
        accountFrom,
        date: new Date(date),
        amount: Number(amount),
        accountTo,
        memo,
        categoryId,
        transferTransactionId,
      } as Transaction)
      .subscribe((transaction: Transaction) => {
        this.transactions.push(transaction);
      });
  }

  updateTransaction(transaction: Transaction): void {
    console.log(transaction);
    this.transactionService.updateTransaction(transaction).subscribe();
  }
  deleteTransaction(transaction: Transaction): void {
    console.log(transaction);
    this.transactions = this.transactions.filter((h) => h !== transaction);
    this.transactionService
      .deleteTransaction(transaction.transactionId)
      .subscribe();
  }

  getCategories(): void {
    this.categoryService
      .getCategories()
      .subscribe((categories) => (this.categories = categories));
  }

  getAccounts(): void {
    this.accountService.getAccounts().subscribe((accounts) => {
      this.accounts = accounts;
      this.accountSelected = accounts[0];
    });
  }

  getAccountNameById(accountId: string) {
    return this.accounts.find((account) => account.accountId == accountId)
      ?.name;
  }
  // getCategorytNameById(categoryId: string): string {
  //   return (
  //     this.categories.find(
  //       (category) => category.categoryId == categoryId
  //     ) as Category
  //   ).name;
  // }

  getCategorytById(categoryId: string): Category {
    return this.categories.find(
      (category) => category.categoryId == categoryId
    ) as Category;
  }

  getCategorytFullName(categoryId: string) {
    let fullName: string = "";

    let category: Category;

    do {
      category = this.getCategorytById(categoryId);

      if (!fullName) {
        fullName = category.name;
      } else {
        fullName = [category.name, fullName].join(" :: ");
      }

      categoryId = category.parentId;
    } while (categoryId);

    return fullName;
  }

  displayedColumns: string[] = [
    "date",
    "accountFrom",
    "amount",
    "accountTo",
    "memo",
    "category",
    "actions",
  ];

  openAddTransactionDialog(): void {
    const dialogRef = this.dialog.open(TransactionAddComponent, {
      width: "500px",
      data: {},
    });

    dialogRef.afterClosed().subscribe((transaction) => {
      console.log("The dialog was closed");
      if (transaction) {
        this.addTransaction(
          transaction.accountFrom,
          transaction.date,
          transaction.amount,
          transaction.accountTo,
          transaction.memo,
          transaction.categoryId,
          transaction.transferTransactionId
        );
      }
    });
  }

  openEditTransactionDialog(transaction: Transaction): void {
    const dialogRef = this.dialog.open(TransactionEditComponent, {
      width: "500px",
      data: transaction,
    });

    dialogRef.afterClosed().subscribe((transaction) => {
      console.log("The dialog was closed");
      if (transaction) {
        this.updateTransaction(transaction);
      }
    });
  }

  openDeleteTransactionDialog(transaction: Transaction): void {
    const dialogRef = this.dialog.open(TransactionDeleteComponent, {
      width: "500px",
      data: transaction,
    });

    dialogRef.afterClosed().subscribe((transaction) => {
      console.log("The dialog was closed");
      if (transaction) {
        this.deleteTransaction(transaction);
      }
    });
  }
}

@Component({
  selector: "transaction-add-dialog",
  templateUrl: "transaction.dialog.html",
})
export class TransactionAddComponent {
  constructor(
    public dialogRef: MatDialogRef<TransactionAddComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Transaction
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

@Component({
  selector: "transaction-edit-dialog",
  templateUrl: "transaction.dialog.html",
})
export class TransactionEditComponent {
  constructor(
    public dialogRef: MatDialogRef<TransactionEditComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Transaction
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

@Component({
  selector: "transaction-delete-dialog",
  templateUrl: "transaction.delete.html",
})
export class TransactionDeleteComponent {
  constructor(
    public dialogRef: MatDialogRef<TransactionDeleteComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Transaction
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

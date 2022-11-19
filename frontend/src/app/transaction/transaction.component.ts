import { Component, OnInit } from "@angular/core";
import { Account } from "../account";
import { AccountService } from "../account.service";
import { Category } from "../category";
import { CategoryService } from "../category.service";
import { Transaction } from "../transaction";
import { TransactionService } from "../transaction.service";

@Component({
  selector: "app-transaction",
  templateUrl: "./transaction.component.html",
  styleUrls: ["./transaction.component.css"],
})
export class TransactionComponent implements OnInit {
  transactions: Transaction[] = [];
  categories: Category[] = [];
  accounts: Account[] = [];

  constructor(
    private transactionService: TransactionService,
    private categoryService: CategoryService,
    private accountService: AccountService
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
    this.accountService
      .getAccounts()
      .subscribe((accounts) => (this.accounts = accounts));
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
}

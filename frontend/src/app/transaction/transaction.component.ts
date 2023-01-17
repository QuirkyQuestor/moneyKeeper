import { Component, OnInit, Inject } from "@angular/core";
import { Account } from "../account";
import { AccountService } from "../account.service";
import { Category } from "../category";
import { CategoryService } from "../category.service";
import { Transaction } from "../transaction";
import { TransactionService } from "../transaction.service";
import * as _ from "lodash-es";
import {
  ModalDismissReasons,
  NgbModal,
  NgbActiveModal,
} from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: "app-transaction",
  templateUrl: "./transaction.component.html",
  styleUrls: ["./transaction.component.css"],
})
export class TransactionComponent implements OnInit {
  transactions: Transaction[] = [];
  categories: Category[] = [];
  accounts: Account[] = [];
  selectedAccount: Account = {
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
    private modalService: NgbModal
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
    // console.debug(transaction);
    this.transactionService.updateTransaction(transaction).subscribe();
  }

  deleteTransaction(transaction: Transaction): void {
    // console.debug(transaction);
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
      this.selectedAccount = accounts[0];
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

  closeResult = "";

  private getDismissReason(reason: any): string {
    if (reason === ModalDismissReasons.ESC) {
      return "by pressing ESC";
    } else if (reason === ModalDismissReasons.BACKDROP_CLICK) {
      return "by clicking on a backdrop";
    } else {
      return `with: ${reason}`;
    }
  }

  openAddTransactionDialog() {
    let transaction: Transaction = {
      transactionId: "",
      accountFrom: "",
      date: new Date(),
      amount: 0,
      accountTo: "",
      memo: "",
      categoryId: "",
      transferTransactionId: "",
    };
    const modalRef = this.modalService.open(TransactionAddComponent);
    modalRef.componentInstance.transaction = transaction;
    modalRef.componentInstance.accounts = this.accounts;
    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        if (result) {
          if (result) {
            this.addTransaction(
              result.accountFrom,
              result.date,
              result.amount,
              result.accountTo,
              result.memo,
              result.categoryId,
              result.transferTransactionId
            );
          }
        }

        this.closeResult = `Closed with: ${JSON.stringify(result)}`;
      },
      (reason) => {
        this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
      }
    );
  }

  openEditTransactionDialog(transaction: Transaction) {
    const modalRef = this.modalService.open(TransactionEditComponent);
    let tr = _.cloneDeep(transaction);
    modalRef.componentInstance.transaction = transaction;
    modalRef.componentInstance.accounts = this.accounts;

    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        if (result && !_.isEqual(result, tr)) {
          this.updateTransaction(result);
        }
        this.closeResult = `Closed with: ${JSON.stringify(result)}`;
      },
      (reason) => {
        this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
      }
    );
  }

  openDeleteTransactionDialog(transaction: Transaction) {
    const modalRef = this.modalService.open(TransactionDeleteComponent);
    modalRef.componentInstance.accountType = transaction;
    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        if (result) {
          this.deleteTransaction(result);
        }

        this.closeResult = `Closed with: ${JSON.stringify(result)}`;
      },
      (reason) => {
        this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
      }
    );
  }
}

@Component({
  styleUrls: ["transaction.dialog.css"],
  selector: "transaction-add-dialog",
  templateUrl: "transaction.dialog.html",
})
export class TransactionAddComponent {
  transaction!: Transaction;
  accounts!: Account[];
  categories!: Category[];
  constructor(public activeModal: NgbActiveModal) {}
}

@Component({
  styleUrls: ["transaction.dialog.css"],
  selector: "transaction-edit-dialog",
  templateUrl: "transaction.dialog.html",
})
export class TransactionEditComponent {
  transaction!: Transaction;
  accounts!: Account[];
  categories!: Category[];
  constructor(public activeModal: NgbActiveModal) {}
}

@Component({
  selector: "transaction-delete-dialog",
  templateUrl: "transaction.delete.html",
})
export class TransactionDeleteComponent {
  transaction!: Transaction;
  constructor(public activeModal: NgbActiveModal) {}
}

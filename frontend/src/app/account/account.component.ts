import { Component, OnInit, Input } from "@angular/core";
import { AccountService } from "../account.service";
import { Account } from "../account";
import { AccountTypeService } from "../account-type.service";
import * as _ from "lodash-es";

import { AccountType } from "../account-type";

import {
  ModalDismissReasons,
  NgbModal,
  NgbActiveModal,
} from "@ng-bootstrap/ng-bootstrap";

export interface DialogData {
  account: Account;
  accountTypes: AccountType[];
}

@Component({
  selector: "app-account",
  templateUrl: "./account.component.html",
  styleUrls: ["./account.component.css"],
})
export class AccountComponent implements OnInit {
  accounts: Account[] = [];
  accountTypes: AccountType[] = [];

  constructor(
    private accountService: AccountService,
    public accountTypeService: AccountTypeService,
    private modalService: NgbModal
  ) {}

  ngOnInit(): void {
    this.getAccountTypes();
    this.getAccounts();
  }

  getAccounts(): void {
    this.accountService
      .getAccounts()
      .subscribe((accounts) => (this.accounts = accounts));
  }

  getAccountTypes(): void {
    this.accountTypeService
      .getAccountTypes()
      .subscribe((types) => (this.accountTypes = types));
  }

  addAccount(
    typeId: string,
    name: string,
    description: string,
    active: boolean
  ): void {
    name = name.trim();
    if (!name) {
      return;
    }
    this.accountService
      .addAccount({ typeId, name, description, active } as Account)
      .subscribe((account: Account) => {
        this.accounts.push(account);
      });
  }

  updateAccount(account: Account): void {
    console.log(account);
    this.accountService.updateAccount(account).subscribe();
  }

  deleteAccount(account: Account): void {
    console.log(account);
    this.accounts = this.accounts.filter((h) => h !== account);
    this.accountService.deleteAccount(account.accountId).subscribe();
  }

  getAccountTypeNameById(typeId: string) {
    return this.accountTypes.find((type) => type.typeId == typeId)?.name;
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

  openAddAccountDialog() {
    let account: Account = {
      accountId: "",
      typeId: this.accountTypes[0].typeId,
      name: "",
      description: "",
      active: true,
    };
    const modalRef = this.modalService.open(AccountAddComponent);
    modalRef.componentInstance.account = account;
    modalRef.componentInstance.accountTypes = this.accountTypes;
    modalRef.result.then(
      (result: Account) => {
        console.log("The dialog was closed");
        if (result) {
          console.log(result);

          this.addAccount(
            result.typeId,
            result.name,
            result.description,
            result.active
          );
        }
        this.closeResult = `Closed with: ${JSON.stringify(result)}`;
      },
      (reason) => {
        this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
      }
    );
  }

  openEditAccountDialog(account: Account) {
    let ac = _.cloneDeep(account);
    const modalRef = this.modalService.open(AccountEditComponent);
    modalRef.componentInstance.account = account;
    modalRef.componentInstance.accountTypes = this.accountTypes;
    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        if (result && !_.isEqual(result, ac)) {
          this.updateAccount(account);
        }

        this.closeResult = `Closed with: ${JSON.stringify(result)}`;
      },
      (reason) => {
        this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
      }
    );
  }

  openDeleteAccountDialog(account: Account) {
    const modalRef = this.modalService.open(AccountDeleteComponent);
    modalRef.componentInstance.account = account;
    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        if (result) {
          this.deleteAccount(account);
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
  styleUrls: ["account.dialog.css"],
  selector: "account-add-dialog",
  templateUrl: "account.dialog.html",
})
export class AccountAddComponent {
  @Input()
  account!: Account;
  accountTypes!: AccountType[];
  constructor(public activeModal: NgbActiveModal) {}
}

@Component({
  styleUrls: ["account.dialog.css"],
  selector: "account-edit-dialog",
  templateUrl: "account.dialog.html",
})
export class AccountEditComponent {
  account!: Account;
  accountTypes!: AccountType[];
  constructor(public activeModal: NgbActiveModal) {}
}

@Component({
  selector: "account-delete-dialog",
  templateUrl: "account.delete.html",
})
export class AccountDeleteComponent {
  account!: Account;
  constructor(public activeModal: NgbActiveModal) {}
}

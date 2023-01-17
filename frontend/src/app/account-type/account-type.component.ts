import { Component, OnInit, Input } from "@angular/core";
import { AccountType } from "../account-type";
import { AccountTypeService } from "../account-type.service";
import * as _ from "lodash-es";
import {
  ModalDismissReasons,
  NgbModal,
  NgbActiveModal,
} from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: "app-account-type",
  templateUrl: "./account-type.component.html",
  styleUrls: ["./account-type.component.css"],
})
export class AccountTypeComponent implements OnInit {
  accountTypes: AccountType[] = [];

  constructor(
    private accountTypeService: AccountTypeService,
    private modalService: NgbModal
  ) {}

  ngOnInit(): void {
    this.getAccountTypes();
  }

  getAccountTypes(): void {
    this.accountTypeService
      .getAccountTypes()
      .subscribe((types) => (this.accountTypes = types));
  }

  addAccountType(name: string, description: string): void {
    name = name.trim();
    if (!name) {
      return;
    }
    this.accountTypeService
      .addAccountType({ name, description } as AccountType)
      .subscribe((accountType: AccountType) => {
        this.accountTypes.push(accountType);
      });
  }

  updateAccountType(accountType: AccountType): void {
    console.log(accountType);
    this.accountTypeService.updateAccountType(accountType).subscribe();
  }

  deleteAccountType(accountType: AccountType): void {
    console.log(accountType);
    this.accountTypes = this.accountTypes.filter((h) => h !== accountType);
    this.accountTypeService.deleteAccountType(accountType.typeId).subscribe();
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

  openEditAccountTypeDialog(accountType: AccountType) {
    const modalRef = this.modalService.open(AccountTypeEditComponent);
    let at = _.cloneDeep(accountType);
    modalRef.componentInstance.accountType = accountType;
    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        if (result && !_.isEqual(result, at)) {
          this.updateAccountType(result);
        }

        this.closeResult = `Closed with: ${JSON.stringify(result)}`;
      },
      (reason) => {
        this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
      }
    );
  }

  openAddAccountTypeDialog() {
    let accountType: AccountType = {
      typeId: "",
      name: "",
      description: "",
    };
    const modalRef = this.modalService.open(AccountTypeAddComponent);
    modalRef.componentInstance.accountType = accountType;
    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        if (result) {
          this.addAccountType(result.name, result.description);
        }

        this.closeResult = `Closed with: ${JSON.stringify(result)}`;
      },
      (reason) => {
        this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
      }
    );
  }

  openDeletedAccountTypeDialog(accountType: AccountType) {
    const modalRef = this.modalService.open(AccountTypeDeleteComponent);
    modalRef.componentInstance.accountType = accountType;
    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        if (result) {
          this.deleteAccountType(accountType);
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
  styleUrls: ["account-type.dialog.css"],
  selector: "account-type-add-dialog",
  templateUrl: "account-type.dialog.html",
})
export class AccountTypeAddComponent {
  accountType!: AccountType;
  constructor(public activeModal: NgbActiveModal) {}
}

@Component({
  styleUrls: ["account-type.dialog.css"],
  selector: "account-type-edit-dialog",
  templateUrl: "account-type.dialog.html",
})
export class AccountTypeEditComponent {
  // @Input()
  accountType!: AccountType;
  constructor(public activeModal: NgbActiveModal) {}
}

@Component({
  selector: "account-type-delete-dialog",
  templateUrl: "account-type.delete.html",
})
export class AccountTypeDeleteComponent {
  // @Input()
  accountType!: AccountType;
  constructor(public activeModal: NgbActiveModal) {}
}

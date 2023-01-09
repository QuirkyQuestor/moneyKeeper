import { Component, OnInit, Inject } from "@angular/core";
import { AccountService } from "../account.service";
import { Account } from "../account";
import { AccountTypeService } from "../account-type.service";

import {
  MatDialog,
  MAT_DIALOG_DATA,
  MatDialogRef,
} from "@angular/material/dialog";
import { AccountType } from "../account-type";

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
    public dialog: MatDialog
  ) {}

  displayedColumns: string[] = [
    "name",
    "account_type",
    "description",
    "active",
    "actions",
  ];

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

  openAddAccountDialog(): void {
    const dialogRef = this.dialog.open(AccountAddComponent, {
      width: "250px",
      data: {
        account: {
          typeId: this.accountTypes[0].typeId,
          active: true,
          accountId: "",
          name: "",
          description: "",
        } as Account,
        accountTypes: this.accountTypes,
      },
    });

    dialogRef.afterClosed().subscribe((data) => {
      console.log("The dialog was closed");
      console.log(JSON.stringify(data));
      if (data.account) {
        this.addAccount(
          data.account.typeId,
          data.account.name,
          data.account.description,
          data.account.active
        );
      }
    });
  }

  openEditAccountDialog(account: Account): void {
    const dialogRef = this.dialog.open(AccountEditComponent, {
      width: "250px",
      data: { account, accountTypes: this.accountTypes },
    });

    dialogRef.afterClosed().subscribe((data) => {
      console.log("The dialog was closed");
      if (data.account) {
        this.updateAccount(data.account);
      }
    });
  }

  openDeleteAccountDialog(account: Account): void {
    const dialogRef = this.dialog.open(AccountDeleteComponent, {
      width: "250px",
      data: account,
    });

    dialogRef.afterClosed().subscribe((account) => {
      console.log("The dialog was closed");
      if (account) {
        this.deleteAccount(account);
      }
    });
  }
}

@Component({
  selector: "account-add-dialog",
  templateUrl: "account.dialog.html",
})
export class AccountAddComponent {
  constructor(
    public dialogRef: MatDialogRef<AccountAddComponent>,
    @Inject(MAT_DIALOG_DATA) public data: DialogData
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

@Component({
  selector: "account-edit-dialog",
  templateUrl: "account.dialog.html",
})
export class AccountEditComponent {
  constructor(
    public dialogRef: MatDialogRef<AccountEditComponent>,
    @Inject(MAT_DIALOG_DATA) public data: DialogData
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

@Component({
  selector: "account-delete-dialog",
  templateUrl: "account.delete.html",
})
export class AccountDeleteComponent {
  constructor(
    public dialogRef: MatDialogRef<AccountDeleteComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Account
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

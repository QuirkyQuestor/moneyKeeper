import { Component, OnInit, Inject } from "@angular/core";
import { AccountType } from "../account-type";
import { AccountTypeService } from "../account-type.service";
import {
  MatDialog,
  MAT_DIALOG_DATA,
  MatDialogRef,
} from "@angular/material/dialog";

@Component({
  selector: "app-account-type",
  templateUrl: "./account-type.component.html",
  styleUrls: ["./account-type.component.css"],
})
export class AccountTypeComponent implements OnInit {
  accountTypes: AccountType[] = [];

  constructor(
    private accountTypeService: AccountTypeService,
    public dialog: MatDialog
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

  displayedColumns: string[] = ["name", "description", "actions"];

  openAddAccountTypeDialog(): void {
    const dialogRef = this.dialog.open(AccountTypeAddComponent, {
      width: "250px",
      data: {},
    });

    dialogRef.afterClosed().subscribe((accountType) => {
      console.log("The dialog was closed");
      if (accountType) {
        this.addAccountType(accountType.name, accountType.description);
      }
    });
  }

  openEditAccountTypeDialog(accountType: AccountType): void {
    const dialogRef = this.dialog.open(AccountTypeEditComponent, {
      width: "250px",
      data: accountType,
    });

    dialogRef.afterClosed().subscribe((accountType) => {
      console.log("The dialog was closed");
      if (accountType) {
        this.updateAccountType(accountType);
      }
    });
  }

  openDeleteAccountTypeDialog(accountType: AccountType): void {
    const dialogRef = this.dialog.open(AccountTypeDeleteComponent, {
      width: "250px",
      data: accountType,
    });

    dialogRef.afterClosed().subscribe((accountType) => {
      console.log("The dialog was closed");
      if (accountType) {
        this.deleteAccountType(accountType);
      }
    });
  }
}

@Component({
  selector: "account-type-add-dialog",
  templateUrl: "account-type.dialog.html",
})
export class AccountTypeAddComponent {
  constructor(
    public dialogRef: MatDialogRef<AccountTypeAddComponent>,
    @Inject(MAT_DIALOG_DATA) public data: AccountType
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

@Component({
  selector: "account-type-edit-dialog",
  templateUrl: "account-type.dialog.html",
})
export class AccountTypeEditComponent {
  constructor(
    public dialogRef: MatDialogRef<AccountTypeEditComponent>,
    @Inject(MAT_DIALOG_DATA) public data: AccountType
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

@Component({
  selector: "account-type-delete-dialog",
  templateUrl: "account-type.delete.html",
})
export class AccountTypeDeleteComponent {
  constructor(
    public dialogRef: MatDialogRef<AccountTypeDeleteComponent>,
    @Inject(MAT_DIALOG_DATA) public data: AccountType
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

import { Component, OnInit } from "@angular/core";
import { AccountType } from "../account-type";
import { AccountTypeService } from "../account-type.service";

@Component({
  selector: "app-account-type",
  templateUrl: "./account-type.component.html",
  styleUrls: ["./account-type.component.css"],
})
export class AccountTypeComponent implements OnInit {
  accountTypes: AccountType[] = [];

  constructor(private accountTypeService: AccountTypeService) {}

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

  deleteAccountType(accountType: AccountType): void {
    console.log(accountType);
    this.accountTypes = this.accountTypes.filter((h) => h !== accountType);
    this.accountTypeService.deleteAccountType(accountType.typeId).subscribe();
  }

  displayedColumns: string[] = ["name", "description", "actions"];
}

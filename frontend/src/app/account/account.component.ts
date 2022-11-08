import { Component, OnInit } from "@angular/core";
import { AccountService } from "../account.service";
import { Account } from "../account";

@Component({
  selector: "app-account",
  templateUrl: "./account.component.html",
  styleUrls: ["./account.component.css"],
})
export class AccountComponent implements OnInit {
  accounts: Account[] = [];

  constructor(private accountService: AccountService) {}

  displayedColumns: string[] = [
    "name",
    "account_type",
    "description",
    "active",
    "actions",
  ];

  ngOnInit(): void {
    this.getAccounts();
  }

  getAccounts(): void {
    this.accountService
      .getAccounts()
      .subscribe((accounts) => (this.accounts = accounts));
  }

  addAccount(name: string, description: string): void {
    name = name.trim();
    if (!name) {
      return;
    }
    this.accountService
      .addAccount({ name, description } as Account)
      .subscribe((account: Account) => {
        this.accounts.push(account);
      });
  }

  deleteAccount(account: Account): void {
    console.log(account);
    this.accounts = this.accounts.filter((h) => h !== account);
    this.accountService.deleteAccount(account.accountId).subscribe();
  }
}

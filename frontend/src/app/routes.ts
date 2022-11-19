import { Routes } from "@angular/router";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { AccountTypeComponent } from "./account-type/account-type.component";
import { AccountComponent } from "./account/account.component";
import { CategoryComponent } from "./category/category.component";
import { TransactionComponent } from "./transaction/transaction.component";

export const routes: Routes = [
  { path: "", title: "Home", redirectTo: "/dashboard", pathMatch: "full" },
  { path: "dashboard", title: "Dashboard", component: DashboardComponent },
  {
    path: "account_type",
    title: "Account Types",
    component: AccountTypeComponent,
  },
  {
    path: "account",
    title: "Accounts",
    component: AccountComponent,
  },
  {
    path: "category",
    title: "Categories",
    component: CategoryComponent,
  },
  {
    path: "transaction",
    title: "Transactions",
    component: TransactionComponent,
  },
];

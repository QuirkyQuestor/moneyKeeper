import { Routes } from "@angular/router";
import { DashboardComponent } from "./dashboard/dashboard.component";
import { AccountTypeComponent } from "./account-type/account-type.component";
import { AccountComponent } from "./account/account.component";
import { CategoryComponent } from "./category/category.component";

export const routes: Routes = [
  { path: "", title: "Home", redirectTo: "/dashboard", pathMatch: "full" },
  { path: "dashboard", title: "Dashboard", component: DashboardComponent },
  {
    path: "account_type",
    title: "AccountType",
    component: AccountTypeComponent,
  },
  {
    path: "account",
    title: "Account",
    component: AccountComponent,
  },
  {
    path: "category",
    title: "Category",
    component: CategoryComponent,
  },
];

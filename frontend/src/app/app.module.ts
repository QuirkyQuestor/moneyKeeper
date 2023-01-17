import { NgModule } from "@angular/core";
import { BrowserModule } from "@angular/platform-browser";
import { HttpClientModule } from "@angular/common/http";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { DashboardComponent } from "./dashboard/dashboard.component";

import {
  AccountTypeComponent,
  AccountTypeAddComponent,
  AccountTypeEditComponent,
  AccountTypeDeleteComponent,
} from "./account-type/account-type.component";
import {
  AccountComponent,
  AccountAddComponent,
  AccountEditComponent,
  AccountDeleteComponent,
} from "./account/account.component";
import {
  CategoryComponent,
  CategoryAddComponent,
  CategoryEditComponent,
  CategoryDeleteComponent,
} from "./category/category.component";
import {
  TransactionComponent,
  TransactionAddComponent,
  TransactionEditComponent,
  TransactionDeleteComponent,
} from "./transaction/transaction.component";

import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { NgbModule } from "@ng-bootstrap/ng-bootstrap";

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,

    AccountTypeComponent,
    AccountTypeAddComponent,
    AccountTypeEditComponent,
    AccountTypeDeleteComponent,

    AccountComponent,
    AccountAddComponent,
    AccountEditComponent,
    AccountDeleteComponent,

    CategoryComponent,
    CategoryAddComponent,
    CategoryEditComponent,
    CategoryDeleteComponent,

    TransactionComponent,
    TransactionAddComponent,
    TransactionEditComponent,
    TransactionDeleteComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    BrowserAnimationsModule,
    FormsModule,
    ReactiveFormsModule,
    NgbModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}

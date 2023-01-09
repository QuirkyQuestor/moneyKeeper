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
import { MatTableModule } from "@angular/material/table";
import { MatDividerModule } from "@angular/material/divider";
import { MatButtonModule } from "@angular/material/button";
import { MatDialogModule } from "@angular/material/dialog";
import { MatFormFieldModule } from "@angular/material/form-field";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { MatInputModule } from "@angular/material/input";
import { MatCheckboxModule } from "@angular/material/checkbox";
import { MatSelectModule } from "@angular/material/select";

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
    MatTableModule,
    MatDividerModule,
    MatButtonModule,
    MatDialogModule,
    MatFormFieldModule,
    FormsModule,
    ReactiveFormsModule,
    // MaterialExampleModule,
    MatInputModule,
    MatCheckboxModule,
    MatSelectModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}

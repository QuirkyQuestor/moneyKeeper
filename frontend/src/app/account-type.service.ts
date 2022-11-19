import { Injectable } from "@angular/core";
import { Observable, of } from "rxjs";
// import { MessageService } from './message.service';
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { catchError, map, tap } from "rxjs/operators";
import { AccountType } from "./account-type";

@Injectable({
  providedIn: "root",
})
export class AccountTypeService {
  constructor(private http: HttpClient) {}

  private accountTypeUrl = "http://localhost:8000/api/account_type"; // URL to web api

  httpOptions = {
    headers: new HttpHeaders({ "Content-Type": "application/json" }),
  };

  /** GET AccountType from the server */
  getAccountTypes(): Observable<AccountType[]> {
    return this.http.get<AccountType[]>(this.accountTypeUrl).pipe(
      tap((_) => console.info("Got AccountTypes")),
      catchError(this.handleError<AccountType[]>("getAccountTypes", []))
    );
  }

  /** POST AccountType from the server */
  addAccountType(accountType: AccountType): Observable<AccountType> {
    return this.http
      .post<AccountType>(this.accountTypeUrl, accountType, this.httpOptions)
      .pipe(
        tap((_) => console.info("Added AccountType")),
        catchError(this.handleError<AccountType>("addAccountTypes"))
      );
  }

  deleteAccountType(account_type_id: string): Observable<AccountType> {
    const url = `${this.accountTypeUrl}/${account_type_id}`;

    return this.http.delete<AccountType>(url, this.httpOptions).pipe(
      tap((_) => console.log(`Deleted AccountType id=${account_type_id}`)),
      catchError(this.handleError<AccountType>("deleteAccountType"))
    );
  }

  updateAccountType(accountType: AccountType): Observable<AccountType> {
    const url = `${this.accountTypeUrl}/${accountType.typeId}`;

    return this.http.put<AccountType>(url, accountType, this.httpOptions).pipe(
      tap((_) => console.log(`Updated AccountType id=${accountType.typeId}`)),
      catchError(this.handleError<AccountType>("updateAccountType"))
    );
  }

  /**
   * Handle Http operation that failed.
   * Let the app continue.
   *
   * @param operation - name of the operation that failed
   * @param result - optional value to return as the observable result
   */
  private handleError<T>(operation = "operation", result?: T) {
    return (error: any): Observable<T> => {
      // TODO: send the error to remote logging infrastructure
      console.error(error); // log to console instead

      // TODO: better job of transforming error for user consumption
      console.log(`${operation} failed: ${error.message}`);

      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
}

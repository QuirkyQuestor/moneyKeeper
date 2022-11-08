import { Injectable } from "@angular/core";
import { Observable, of } from "rxjs";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { catchError, map, tap } from "rxjs/operators";
import { Account } from "./account";

@Injectable({
  providedIn: "root",
})
export class AccountService {
  constructor(private http: HttpClient) {}

  private AccountUrl = "http://localhost:8000/api/account"; // URL to web api

  httpOptions = {
    headers: new HttpHeaders({ "Content-Type": "application/json" }),
  };

  /** GET Account from the server */
  getAccounts(): Observable<Account[]> {
    return this.http.get<Account[]>(this.AccountUrl).pipe(
      tap((_) => console.info("Got Accounts")),
      catchError(this.handleError<Account[]>("getAccounts", []))
    );
  }

  /** POST Account from the server */
  addAccount(account: Account): Observable<Account> {
    return this.http
      .post<Account>(this.AccountUrl, account, this.httpOptions)
      .pipe(
        tap((_) => console.info("Added Account")),
        catchError(this.handleError<Account>("addAccounts"))
      );
  }

  deleteAccount(account_id: string): Observable<Account> {
    const url = `${this.AccountUrl}/${account_id}`;

    return this.http.delete<Account>(url, this.httpOptions).pipe(
      tap((_) => console.log(`Deleted Account id=${account_id}`)),
      catchError(this.handleError<Account>("deleteAccount"))
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

import { Injectable } from "@angular/core";
import { Transaction } from "./transaction";
import { Observable, of } from "rxjs";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { catchError, map, tap } from "rxjs/operators";

@Injectable({
  providedIn: "root",
})
export class TransactionService {
  constructor(private http: HttpClient) {}

  private TransactionUrl = "http://localhost:8000/api/transaction"; // URL to web api

  httpOptions = {
    headers: new HttpHeaders({ "Content-Type": "application/json" }),
  };

  /** GET Transaction from the server */
  getTransactions(): Observable<Transaction[]> {
    return this.http.get<Transaction[]>(this.TransactionUrl).pipe(
      tap((_) => console.info("Got Transactions")),
      catchError(this.handleError<Transaction[]>("getTransactions", []))
    );
  }

  /** POST Transaction from the server */
  addTransaction(transaction: Transaction): Observable<Transaction> {
    console.info("Added Transaction", transaction);
    return this.http
      .post<Transaction>(this.TransactionUrl, transaction, this.httpOptions)
      .pipe(
        tap((_) => console.info("Added Transaction")),
        catchError(this.handleError<Transaction>("addTransaction"))
      );
  }

  updateTransaction(transaction: Transaction): Observable<Transaction> {
    console.info("Update Transaction", transaction);
    const url = `${this.TransactionUrl}/${transaction.transactionId}`;

    return this.http.put<Transaction>(url, transaction, this.httpOptions).pipe(
      tap((_) => console.info("Updated Transaction")),
      catchError(this.handleError<Transaction>("updateTransaction"))
    );
  }

  deleteTransaction(transaction_id: string): Observable<Transaction> {
    const url = `${this.TransactionUrl}/${transaction_id}`;

    return this.http.delete<Transaction>(url, this.httpOptions).pipe(
      tap((_) => console.log(`Deleted Transaction id=${transaction_id}`)),
      catchError(this.handleError<Transaction>("deleteTransaction"))
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

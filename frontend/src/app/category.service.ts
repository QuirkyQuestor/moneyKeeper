import { Injectable } from "@angular/core";
import { Category } from "./category";
import { Observable, of } from "rxjs";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { catchError, map, tap } from "rxjs/operators";

@Injectable({
  providedIn: "root",
})
export class CategoryService {
  constructor(private http: HttpClient) {}

  private CategoryUrl = "http://localhost:8000/api/category"; // URL to web api

  httpOptions = {
    headers: new HttpHeaders({ "Content-Type": "application/json" }),
  };

  /** GET Categories from the server */
  getCategories(): Observable<Category[]> {
    return this.http.get<Category[]>(this.CategoryUrl).pipe(
      tap((_) => console.info("Got Categories")),
      catchError(this.handleError<Category[]>("getCategories", []))
    );
  }

  /** POST Category from the server */
  addCategory(category: Category): Observable<Category> {
    return this.http
      .post<Category>(this.CategoryUrl, category, this.httpOptions)
      .pipe(
        tap((_) => console.info("Added Category")),
        catchError(this.handleError<Category>("addCategory"))
      );
  }

  deleteCategory(category_id: string): Observable<Category> {
    const url = `${this.CategoryUrl}/${category_id}`;

    return this.http.delete<Category>(url, this.httpOptions).pipe(
      tap((_) => console.log(`Deleted Category id=${category_id}`)),
      catchError(this.handleError<Category>("deleteCategory"))
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

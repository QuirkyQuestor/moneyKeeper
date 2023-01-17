import { Component, OnInit, Inject } from "@angular/core";
import { Category } from "../category";
import { CategoryService } from "../category.service";
import * as _ from "lodash-es";
import {
  ModalDismissReasons,
  NgbModal,
  NgbActiveModal,
} from "@ng-bootstrap/ng-bootstrap";

@Component({
  selector: "app-category",
  templateUrl: "./category.component.html",
  styleUrls: ["./category.component.css"],
})
export class CategoryComponent implements OnInit {
  categories: Category[] = [];

  constructor(
    private categoryService: CategoryService,
    private modalService: NgbModal
  ) {}

  ngOnInit(): void {
    this.getCategories();
  }

  getCategories(): void {
    this.categoryService
      .getCategories()
      .subscribe((categories) => (this.categories = categories));
  }

  addCategory(
    name: string,
    parentId: string,
    description: string,
    expence: boolean
  ): void {
    name = name.trim();
    if (!name) {
      return;
    }
    this.categoryService
      .addCategory({ name, parentId, description, expence } as Category)
      .subscribe((category: Category) => {
        this.categories.push(category);
      });
  }

  deleteCategory(category: Category): void {
    console.log(category);
    this.categories = this.categories.filter((h) => h !== category);
    this.categoryService.deleteCategory(category.categoryId).subscribe();
  }

  updateCategory(category: Category): void {
    console.log(category);
    this.categoryService.updateCategory(category).subscribe();
  }

  getCategorytById(categoryId: string): Category {
    return this.categories.find(
      (category) => category.categoryId == categoryId
    ) as Category;
  }

  getCategorytParentName(category: Category) {
    let parentName: string = "";

    while (category.parentId) {
      let parent = this.getCategorytById(category.parentId);

      if (!parentName) {
        parentName = parent.name;
      } else {
        parentName = [parent.name, parentName].join(" :: ");
      }

      category = parent;
    }

    return parentName;
  }

  closeResult = "";

  private getDismissReason(reason: any): string {
    if (reason === ModalDismissReasons.ESC) {
      return "by pressing ESC";
    } else if (reason === ModalDismissReasons.BACKDROP_CLICK) {
      return "by clicking on a backdrop";
    } else {
      return `with: ${reason}`;
    }
  }

  openAddCategoryDialog() {
    const modalRef = this.modalService.open(CategoryAddComponent);
    modalRef.componentInstance.categories = this.categories;

    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        console.log(JSON.stringify(result));
        if (result) {
          this.addCategory(
            result.name,
            result.parentId,
            result.description,
            result.expence
          );
        }
        this.closeResult = `Closed with: ${JSON.stringify(result)}`;
      },
      (reason) => {
        this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
      }
    );
  }

  openEditCategoryDialog(category: Category) {
    let ct = _.cloneDeep(category);

    const modalRef = this.modalService.open(CategoryEditComponent);
    modalRef.componentInstance.category = category;
    modalRef.componentInstance.categories = this.categories;

    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        if (result && !_.isEqual(result, ct)) {
          this.updateCategory(result);
        }
        this.closeResult = `Closed with: ${JSON.stringify(result)}`;
      },
      (reason) => {
        this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
      }
    );
  }

  openDeleteCategoryDialog(category: Category) {
    const modalRef = this.modalService.open(CategoryDeleteComponent);
    modalRef.componentInstance.category = category;
    modalRef.result.then(
      (result) => {
        console.log("The dialog was closed");
        if (result) {
          this.deleteCategory(result);
        }
        this.closeResult = `Closed with: ${JSON.stringify(result)}`;
      },
      (reason) => {
        this.closeResult = `Dismissed ${this.getDismissReason(reason)}`;
      }
    );
  }
}

@Component({
  selector: "category-add-dialog",
  templateUrl: "category.dialog.html",
})
export class CategoryAddComponent {
  // @Input()
  category: Category = {
    categoryId: "",
    parentId: "",
    name: "",
    description: "",
    expence: true,
  };
  categories!: Category[];
  constructor(public activeModal: NgbActiveModal) {}

  getCategorytParentName(category: Category) {
    let parentName: string = "";

    while (category.parentId) {
      let parent = this.getCategorytById(category.parentId);

      if (!parentName) {
        parentName = parent.name;
      } else {
        parentName = [parent.name, parentName].join(" :: ");
      }

      category = parent;
    }

    return parentName;
  }

  getCategorytById(categoryId: string): Category {
    return this.categories.find(
      (category) => category.categoryId == categoryId
    ) as Category;
  }
}

@Component({
  selector: "category-edit-dialog",
  templateUrl: "category.dialog.html",
})
export class CategoryEditComponent {
  category!: Category;
  categories!: Category[];
  constructor(public activeModal: NgbActiveModal) {}

  getCategorytParentName(category: Category) {
    let parentName: string = "";

    while (category.parentId) {
      let parent = this.getCategorytById(category.parentId);

      if (!parentName) {
        parentName = parent.name;
      } else {
        parentName = [parent.name, parentName].join(" :: ");
      }

      category = parent;
    }

    return parentName;
  }

  getCategorytById(categoryId: string): Category {
    return this.categories.find(
      (category) => category.categoryId == categoryId
    ) as Category;
  }
}

@Component({
  selector: "category-delete-dialog",
  templateUrl: "category.delete.html",
})
export class CategoryDeleteComponent {
  category!: Category;
  constructor(public activeModal: NgbActiveModal) {}
}
